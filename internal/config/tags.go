package config

import (
	"fmt"
	"slices"
	"strings"
	"unicode"

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
		_, after, _ := strings.Cut(trimmed, ":")
		content := after
		parts := strings.FieldsFunc(content, func(r rune) bool {
			return r == ',' || r == ' ' || r == '\t' || r == '\n' || r == '\r'
		})
		raw = append(raw, parts...)
	} else {
		fields := strings.FieldsSeq(trimmed)
		for field := range fields {
			if strings.HasPrefix(field, "#") {
				raw = append(raw, field)
			}
		}
	}

	var tags []string
	for _, r := range raw {
		clean := cleanTag(r)
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

// UpdateASTHostTags updates or inserts a tag comment node inside an AST Host block
// and prunes any duplicate or stale tag comment lines.
func UpdateASTHostTags(astHost *ssh_config.Host, tags []string) error {
	if astHost == nil {
		return fmt.Errorf("astHost cannot be nil")
	}

	var sanitized []string
	for _, t := range tags {
		clean := cleanTag(t)
		if clean != "" && !slices.Contains(sanitized, clean) {
			sanitized = append(sanitized, clean)
		}
	}

	firstTagIdx := -1
	var filteredNodes []ssh_config.Node
	for i, node := range astHost.Nodes {
		if empty, ok := node.(*ssh_config.Empty); ok {
			if len(ExtractTagsFromComment(empty.String())) > 0 {
				if firstTagIdx == -1 {
					firstTagIdx = i
				}
				continue
			}
		}
		filteredNodes = append(filteredNodes, node)
	}

	if len(sanitized) == 0 {
		astHost.Nodes = filteredNodes
		return nil
	}

	newNode, err := createTagCommentNode(sanitized)
	if err != nil {
		return err
	}

	if firstTagIdx != -1 && firstTagIdx <= len(filteredNodes) {
		filteredNodes = slices.Insert(filteredNodes, firstTagIdx, newNode)
	} else {
		filteredNodes = append([]ssh_config.Node{newNode}, filteredNodes...)
	}

	astHost.Nodes = filteredNodes
	return nil
}

func cleanTag(s string) string {
	s = strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' || r == '\x00' || !unicode.IsPrint(r) {
			return -1
		}
		return r
	}, s)
	s = strings.TrimSpace(s)
	s = strings.TrimLeft(s, "#")
	return strings.ToLower(strings.TrimSpace(s))
}

func createTagCommentNode(tags []string) (ssh_config.Node, error) {
	line := "    # tags: " + strings.Join(tags, ", ") + "\n"
	decoded, err := ssh_config.Decode(strings.NewReader(line))
	if err != nil || len(decoded.Hosts) == 0 || len(decoded.Hosts[0].Nodes) == 0 {
		return nil, fmt.Errorf("failed to create tag comment AST node: %w", err)
	}
	return decoded.Hosts[0].Nodes[0], nil
}

const serviceMarker = "tusshi: service"

// ExtractServiceMarker returns true when any comment node inside the host block
// contains the "tusshi: service" marker (case-insensitive).
func ExtractServiceMarker(nodes []ssh_config.Node) bool {
	for _, node := range nodes {
		if empty, ok := node.(*ssh_config.Empty); ok {
			line := strings.ToLower(strings.TrimSpace(empty.String()))
			if strings.Contains(line, serviceMarker) {
				return true
			}
		}
	}
	return false
}

// WriteServiceMarker prepends a "# tusshi: service" comment node to the host block
// so the entry is hidden from the tuSSHi connection list on next load.
func WriteServiceMarker(astHost *ssh_config.Host) error {
	if astHost == nil {
		return fmt.Errorf("astHost cannot be nil")
	}

	decoded, err := ssh_config.Decode(strings.NewReader("    # " + serviceMarker + "\n"))
	if err != nil || len(decoded.Hosts) == 0 || len(decoded.Hosts[0].Nodes) == 0 {
		return fmt.Errorf("failed to create service marker AST node: %w", err)
	}

	node := decoded.Hosts[0].Nodes[0]
	astHost.Nodes = append([]ssh_config.Node{node}, astHost.Nodes...)
	return nil
}
