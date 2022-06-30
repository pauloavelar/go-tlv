package tlv

// Nodes shorthand for []Node.
type Nodes []Node

// GetByTag returns nodes that match the tag.
func (ns Nodes) GetByTag(tag Tag) Nodes {
	var res Nodes

	for i := range ns {
		if ns[i].Tag == tag {
			res = append(res, ns[i])
		}
	}

	return res
}

// GetFirstByTag returns nodes that match the tag.
func (ns Nodes) GetFirstByTag(tag Tag) (res Node, ok bool) {
	for i := range ns {
		if ns[i].Tag == tag {
			return ns[i], true
		}
	}

	return res, false
}

// HasTag returns if a tag is present in the nodes.
func (ns Nodes) HasTag(tag Tag) bool {
	for i := range ns {
		if ns[i].Tag == tag {
			return true
		}
	}
	return false
}
