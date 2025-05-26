package compiler

// func WalkComponent(comp *Component, visit func(node *Node)) {
// 	for _, node := range comp.Body {
// 		WalkNode(node, visit)
// 	}
// }

// func WalkNode(n *Node, visit func(node *Node)) {
// 	if n == nil {
// 		return
// 	}

// 	visit(n)

// 	if n.Element != nil {
// 		for _, child := range n.Element.Children {
// 			WalkNode(child, visit)
// 		}
// 	}
// }
