package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/6l20/zama-test/common/log"
)

type MerkleNode struct {
    Hash  string		`json:"hash"`
    Left  *MerkleNode	`json:"left"`
    Right *MerkleNode	`json:"right"`
	Parent *MerkleNode	`json:"parent"`
	Sibling *MerkleNode	`json:"sibling"`
	IsLeft bool			`json:"isLeft"`
}

func (n *MerkleNode) IsLeaf() bool {
	return n.Left == nil && n.Right == nil
}

type Proof struct {
	Proof []*ProofStep		`json:"hash"`
}

type ProofStep struct {
	Hash string			`json:"hash"`
	IsLeft bool			`json:"isLeft"`
}

// A MerkleTree contains the root of the tree.
type MerkleTree struct {
    Root *MerkleNode	`json:"root"`
	Nodes []*MerkleNode	`json:"nodes"`
}

type MerkleManager struct {
	MerkleTree *MerkleTree
	logger log.Logger
	filePath string
	level int
}

func NewMerkleManager(filePath string, logger log.Logger) *MerkleManager {
	return &MerkleManager{
		logger: logger,
		filePath: filePath,
	}
}

// CalculateHash takes data and returns the SHA256 hash as a string.
func (m *MerkleManager) calculateHash(data []byte) string {
    hash := sha256.Sum256(data)
    return hex.EncodeToString(hash[:])
}

// NewMerkleNode creates a new MerkleTree node.
func (m *MerkleManager) NewMerkleNode(left, right *MerkleNode, blockData []byte) *MerkleNode {
    mNode := MerkleNode{}

    if left == nil && right == nil {
        mNode.Hash = m.calculateHash(blockData)
		return &mNode
    }
	// Otherwise, hash the concatenation of the left and right children's hashes.
	prevHashes := []byte(left.Hash + right.Hash)
	mNode.Hash = m.calculateHash(prevHashes)
	mNode.Left = left
	mNode.Right = right
	
	return &mNode
}

func (m *MerkleManager) NewMerkleNodeFromBlock( blockData []byte) *MerkleNode {
	mNode := MerkleNode{}
	mNode.Hash = m.calculateHash(blockData)
	return &mNode
}


// PopulateTree creates a new Merkle Tree from a list of data blocks.
func (m *MerkleManager) PopulateTree(  blocks [][]byte) *MerkleNode {

	if m.MerkleTree == nil {
		m.MerkleTree = &MerkleTree{}
	}

	// Base case: if there is only one block, return a new MerkleNode.
	var nodes []*MerkleNode
	for _, b := range blocks {
		m.logger.Debug("PopulateTree", "block", string(b))
		node := m.NewMerkleNode( nil,nil,b)	
		nodes = append(nodes, node)
		m.MerkleTree.Nodes = append(m.MerkleTree.Nodes, node)
	}

	
	
	
    // Recursive case: split the blocks and build the sub-trees

    for len( nodes ) > 1 {
		var upperLayer []*MerkleNode
		for i := 0; i < len(nodes); i += 2 {
			l := nodes[i]
			l.IsLeft = true
			r := l
			if i+1 < len(nodes) {
				r = nodes[i+1]
			}
			parent := m.NewMerkleNode( l, r, nil)
			upperLayer = append(upperLayer, parent)
			
			l.Parent = parent
			r.Parent = parent
			l.Sibling = r
			r.Sibling = l
		}
		nodes = upperLayer
	}
	return nodes[0]

}

func nodesAsBlocks(nodes []*MerkleNode) [][]byte {
    var blocks [][]byte
    for _, node := range nodes {
        blocks = append(blocks, []byte(node.Hash))
    }
    return blocks
}

// BuildMerkleTreeFromFS walks the file system from the root directory and builds a Merkle Tree.
func (m *MerkleManager) BuildMerkleTreeFromFS(rootDir string) (*MerkleTree, error) {
    var blockData [][]byte

    // Walk the file system from the root directory.
    err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // Skip the root directory itself.
        if path == rootDir {
            return nil
        }

        // Read file and get data for leaf nodes.
        if !info.IsDir() {
            data, err := os.ReadFile(path)
            if err != nil {
                return err
            }

            blockData = append(blockData, data)
        }

        return nil
    })

    if err != nil {
        return nil, err
    }

	rootNode := m.PopulateTree( blockData)

	
	m.MerkleTree.Root = rootNode

    return m.MerkleTree, nil
}


// StoreMerkleRoot writes the Merkle root to a local file.
func (m *MerkleManager) StoreMerkleRoot() error {
	if m.MerkleTree == nil {
		return fmt.Errorf("MerkleTree is nil")
	}

	m.logger.Debug("StoreMerkleRoot","hash", m.MerkleTree.Root.Hash)
	
    // Convert the Merkle root to a byte slice for writing to a file.
    data := []byte(m.MerkleTree.Root.Hash)
    // Write the data to the file, with file permissions set to read and write.
    return os.WriteFile(m.filePath, data, 0644)
}

// LoadMerkleRoot reads the Merkle root from a local file.
func (m *MerkleManager) LoadMerkleRoot() (string, error) {
    // Read the data back from the file.
    data, err := os.ReadFile(m.filePath)
    if err != nil {
        return "", err
    }
    // Convert the byte slice back to a string.
    return string(data), nil
}

func (m *MerkleManager) GenerateProof(leafIndex int) *Proof {
	path := []*ProofStep{}

	leavesCounter := 0
	var node *MerkleNode

	for _, n := range m.MerkleTree.Nodes {
		if (leavesCounter == leafIndex) && n.IsLeaf(){
			node = n
			break
		}
		leavesCounter++
	}

	for node != nil  {
		if node.Sibling == nil {
			break
		}
		path = append(path, &ProofStep{ Hash: node.Sibling.Hash, IsLeft: node.IsLeft })
		
		node = node.Parent
	}

	proof := &Proof{
		Proof: path,
	}

	return proof
}

func (m *MerkleManager) VerifyProof(leafHash string, proof Proof, rootHash string) bool {
	// Concatenate the leaf hash with each proof hash and hash them together.
	calculatedHash := leafHash
	for _, step := range proof.Proof {
		m.logger.Debug("VerifyProof", "calculatedHash", calculatedHash, "step", step)
		if step.IsLeft {
			calculatedHash = m.calculateHash([]byte(calculatedHash + step.Hash))
		} else {
			calculatedHash = m.calculateHash([]byte(step.Hash + calculatedHash))
		}
	}
	m.logger.Debug("VerifyProof end", "calculatedHash", calculatedHash, "rootHash", rootHash)
	// The final hash should match the root hash.
	return calculatedHash == rootHash
}
