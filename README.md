# zama-test
zama-test project (Merkle tree)

# Description

This project consists of both a client and a server.

Client is able to store files by uploading them to the server, and keep a merkle root as an integrity checker of the uploaded files.

Server stores the files and also enables file download as well as merkle proofs for integrity checking.

The project may be run using docker-compose.

Try it yourself with :

```
docker-compose up --build
```

The test/data folder contains files to run tests with.

# Approach of the problem

## Computing the root Hash

1. Divide into data blocks

File contents is the data block ( each file gives one data block )

2. Hash each data block

sha256 Hashing algo is used ( this could made configurable )

3. Create Leaf nodes

First step of the process is to create one node with each Hashed data block.

N = H(dataBlock)

Each of these nodes will be the leaves of the tree.

4. Build the tree

Pair the leaf nodes and compute the hash of their concatenation to form their parent node. Continue this process level by level until you reach a single hash at the top, this last value is the Merkle root.

## Generate proofs

For a given file (corresponding leaf node) each Merkle proof is computed by gathering hashes of pairing nodes by ascending the tree to the top. In the code these pairing nodes are called "siblings".

1. Identify the Leaf Node: Locate the leaf node that corresponds to the data block for which you're proving inclusion.

2. Collect Sibling Nodes: Starting from the leaf node, collect the sibling node at each level of the tree. These siblings are necessary to compute the parent node's hash at each level.

3. Ascend the Tree: Move up the tree, collecting sibling nodes until you reach the root.

4. Create the Proof: The Merkle proof is the ordered list of sibling nodes collected. 

## Verify

In order to verify, one must hold the root hash and the hash of the verified data block.

Verification is computed by reconstructing the targeted branch until the top (root) using the proof data and the known file hash.

It is also important to notice that hashes computation be made in the correct order, Hashes concatenation must follow the right path.

```
computedRoot == storedRoot // true if file is same
```

# References

https://www.wikiwand.com/en/Merkle_tree
