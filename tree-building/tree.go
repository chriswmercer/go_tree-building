//Package tree provides types and functions to convert unsorted records into tree-ified nodes
package tree

import (
	"errors"
	"sort"
)

//Record is the record to tree-ify
type Record struct {
	ID     int
	Parent int
}

//Node is the resultant object
type Node struct {
	ID       int
	Children []*Node
}

func sortRecords(r []Record) {
	sort.Slice(r[:], func(i, j int) bool {
		return r[i].ID < r[j].ID
	})
}

func (n Node) sortChildren() {
	sort.Slice(n.Children[:], func(i, j int) bool {
		return n.Children[i].ID < n.Children[j].ID
	})
}

func where(r *[]Record, test func(Record) bool) (ret *[]Record, remainder *[]Record) {
	ret = new([]Record)
	remainder = new([]Record)

	for _, s := range *r {
		if test(s) {
			*ret = append(*ret, s)
		} else {
			*remainder = append(*remainder, s)
		}
	}
	return
}

func (r Record) hasDuplicatesIn(other []Record) bool {
	rID := r.ID
	var sum int
	for _, elem := range other {
		if elem.ID == rID {
			sum++
		}
	}
	return sum > 1
}

func newNode(r Record) *Node {
	n := new(Node)
	n.ID = r.ID
	return n
}

func marshalChildren(parent *Node, potentialChildren *[]Record) (remainderChildren *[]Record, err error) {
	potentialChildren, remainderChildren = where(potentialChildren, func(r Record) bool { return r.Parent == parent.ID })
	for _, child := range *potentialChildren {
		if child.hasDuplicatesIn(*potentialChildren) {
			return nil, errors.New("duplicates detected")
		}
		if parent.ID > child.ID {
			return nil, errors.New("invalid parent id detected")
		}
		if child.Parent > child.ID {
			return nil, errors.New("invalid parent id detected; pid > cid")
		}
		nodeToAdd := newNode(child)
		(*parent).Children = append((*parent).Children, nodeToAdd)
		if len(*remainderChildren) > 0 {
			remainderChildren, err = marshalChildren(nodeToAdd, remainderChildren)
			if err != nil {
				return nil, err
			}
		}
	}
	parent.sortChildren()
	return remainderChildren, nil
}

//Build will build a tree structure of nodes from a given slice of records
func Build(r []Record) (*Node, error) {
	if len(r) == 0 {
		return nil, nil
	}
	sortRecords(r)
	if r[len(r)-1].ID != len(r)-1 {
		return nil, errors.New("non contiguous ids detected")
	}
	parentRecords, remainder := where(&r, func(r Record) bool { return r.ID == r.Parent })
	if len(*parentRecords) != 1 {
		return nil, errors.New("records must have exactly 1 parent Record where ID == Parent")
	}
	parentRecord := newNode((*parentRecords)[0])
	if parentRecord.ID != 0 {
		return nil, errors.New("records must have exactly 1 parent Record where ID == 0")
	}
	parentRecord.sortChildren()
	finalRemainder, err := marshalChildren(parentRecord, remainder)

	if finalRemainder != nil && len(*finalRemainder) > 0 {
		return nil, errors.New("detected orphans")
	}
	return parentRecord, err
}