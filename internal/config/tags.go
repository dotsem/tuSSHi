package config

import (
	"fmt"
	"slices"
	"strings"

	"github.com/kevinburke/ssh_config"
)

// ExtractTagsFromComment parses tag values from a single comment line.
// Supports `# tags: tag1, tag2` and `# #hashtag` formats.
func ExtractTagsFromComment(line string) []string {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "#") {
		return nil
	}
	trimmed = strings.TrimLeft(trimmed, "#")
	trimmed = strings.TrimSpace(trimmed)
	if trimmed == "" {
		return nil
	}

	var raw []string
	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "tags:") || strings.HasPrefix(lower, "tag:") {
		idx := strings.Index(trimmed, ":")
		content := trimmed[idx+1:]
		parts := strings.FieldsFunc(content, func(r rune) bool {
			return r == ',' || r == ' ' || r == '\t'
		})
		raw = append(raw, parts...)
	} else {
		fields := strings.Fields(trimmed)
		for _, field := range fields {
			if strings.HasPrefix(field, "#") {
				raw = append(raw, strings.TrimPrefix(field, "#"))
			}
		}
	}

	var tags []string
	for _, r := range raw {
		clean := strings.ToLower(strings.TrimSpace(strings.TrimPrefix(r, "#")))
		if clean != "" && !slices.Contains(tags, clean) {
			tags = append(tags, clean)
		}
	}
	return tags
}

// ExtractTagsFromNodes inspects AST nodes inside a Host block and collects all unique tags.
func ExtractTagsFromNodes(nodes []ssh_config.Node) []string {
	var tags []string
	for _, node := range nodes {
		if empty, ok := node.(*ssh_config.Empty); ok {
			line := empty.String()
			for _, t := range ExtractTagsFromComment(line) {
				if !slices.Contains(tags, t) {
					tags = append(tags, t)
				}
			}
		}
	}
	return tags
}

// UpdateASTHostTags updates or inserts a tag comment node inside an AST Host block.
func UpdateASTHostTags(astHost *ssh_config.Host, tags []string) error {
	if astHost == nil {
		return fmt.Errorf("astHost cannot be nil")
	}

	tagIdx := -1
	for i, node := range astHost.Nodes {
		if empty, ok := node.(*ssh_config.Empty); ok {
			parsed := ExtractTagsFromComment(empty.String())
			if len(parsed) > 0 {
				tagIdx = i
				break
			}
		}
	}

	if len(tags) == 0 {
		if tagIdx != -1 {
			astHost.Nodes = append(astHost.Nodes[:tagIdx], astHost.Nodes[tagIdx+1:]...)
		}
		return nil
	}

	newNode, err := createTagCommentNode(tags)
	if err != nil {
		return err
	}

	if tagIdx != -1 {
		astHost.Nodes[tagIdx] = newNode
	} else {
		astHost.Nodes = append([]ssh_config.Node{newNode}, astHost.Nodes...)
	}
	return nil
}

func createTagCommentNode(tags []string) (ssh_config.Node, error) {
	line := "    # tags: " + strings.Join(tags, ", ") + "\n"
	decoded, err := ssh_config.Decode(strings.NewReader(line))
	if err != nil || len(decoded.Hosts) == 0 || len(decoded.Hosts[0].Nodes) == 0 {
		return nil, fmt.Errorf("failed to create tag comment AST node: %w", err)
	}
	return decoded.Hosts[0].Nodes[0], nil
}
