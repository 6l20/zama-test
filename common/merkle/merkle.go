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
	level int
}

func NewMerkleManager(logger log.Logger) *MerkleManager {
	return &MerkleManager{
		logger: logger,
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

func (m *MerkleManager) GenerateProof(leafIndex int) (*Proof, error) {
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

	if node == nil {
		return nil, fmt.Errorf("Leaf not found")
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

	return proof, nil
}

func (m *MerkleManager) VerifyProof(leafHash string, proof Proof, rootHash string) bool {
	// Check leaf hash.
	if leafHash == "" {
		m.logger.Warn("VerifyProof", "leafHash", "empty")
		return false
	}
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
