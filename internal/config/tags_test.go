package config

import (
	"strings"
	"testing"

	"github.com/kevinburke/ssh_config"
	"github.com/stretchr/testify/assert"
)

func TestExtractTagsFromComment(t *testing.T) {
	t.Run("keyed comment with commas", func(t *testing.T) {
		line := "# tags: production, database, aws"
		tags := ExtractTagsFromComment(line)
		assert.Equal(t, []string{"production", "database", "aws"}, tags)
	})

	t.Run("keyed comment with spaces and mixed case", func(t *testing.T) {
		line := "   # TAGS: Prod Web AWS  "
		tags := ExtractTagsFromComment(line)
		assert.Equal(t, []string{"prod", "web", "aws"}, tags)
	})

	t.Run("hashtag comment", func(t *testing.T) {
		line := "# #prod #database #aws"
		tags := ExtractTagsFromComment(line)
		assert.Equal(t, []string{"prod", "database", "aws"}, tags)
	})

	t.Run("deduplicates tags", func(t *testing.T) {
		line := "# tags: prod, database, PROD, aws, database"
		tags := ExtractTagsFromComment(line)
		assert.Equal(t, []string{"prod", "database", "aws"}, tags)
	})

	t.Run("non tag comment returns nil", func(t *testing.T) {
		line := "# Standard comment describing host"
		tags := ExtractTagsFromComment(line)
		assert.Nil(t, tags)
	})

	t.Run("non comment line returns nil", func(t *testing.T) {
		line := "HostName 10.0.0.1"
		tags := ExtractTagsFromComment(line)
		assert.Nil(t, tags)
	})
}

func TestExtractTagsFromNodesAndManager(t *testing.T) {
	configContent := `
Host prod-db
    # tags: production, database
    HostName 10.0.0.15
    User postgres

Host staging-web
    # #staging #frontend
    HostName 10.0.0.20
    User deploy
`
	cfg, err := ssh_config.Decode(strings.NewReader(configContent))
	assert.NoError(t, err)

	// ssh_config injects an implicit default host at index 0; explicit hosts start at index 1
	tagsDB := ExtractTagsFromNodes(cfg.Hosts[1].Nodes)
	assert.Equal(t, []string{"production", "database"}, tagsDB)

	tagsWeb := ExtractTagsFromNodes(cfg.Hosts[2].Nodes)
	assert.Equal(t, []string{"staging", "frontend"}, tagsWeb)
}

func TestUpdateASTHostTags(t *testing.T) {
	t.Run("insert new tag comment when none exists", func(t *testing.T) {
		content := "Host test-host\n    HostName 127.0.0.1\n"
		cfg, err := ssh_config.Decode(strings.NewReader(content))
		assert.NoError(t, err)

		err = UpdateASTHostTags(cfg.Hosts[1], []string{"web", "prod"})
		assert.NoError(t, err)

		tags := ExtractTagsFromNodes(cfg.Hosts[1].Nodes)
		assert.Equal(t, []string{"web", "prod"}, tags)
	})

	t.Run("update existing tag comment in place", func(t *testing.T) {
		content := "Host test-host\n    # tags: old1, old2\n    HostName 127.0.0.1\n"
		cfg, err := ssh_config.Decode(strings.NewReader(content))
		assert.NoError(t, err)

		err = UpdateASTHostTags(cfg.Hosts[1], []string{"new1", "new2"})
		assert.NoError(t, err)

		tags := ExtractTagsFromNodes(cfg.Hosts[1].Nodes)
		assert.Equal(t, []string{"new1", "new2"}, tags)
	})

	t.Run("clear tags removes tag node", func(t *testing.T) {
		content := "Host test-host\n    # tags: old1, old2\n    HostName 127.0.0.1\n"
		cfg, err := ssh_config.Decode(strings.NewReader(content))
		assert.NoError(t, err)

		err = UpdateASTHostTags(cfg.Hosts[1], nil)
		assert.NoError(t, err)

		tags := ExtractTagsFromNodes(cfg.Hosts[1].Nodes)
		assert.Empty(t, tags)
	})

	t.Run("sanitizes newlines and control characters", func(t *testing.T) {
		line := "# tags: prod\nHost evil.com\n    User root, aws"
		tags := ExtractTagsFromComment(line)
		assert.Equal(t, []string{"prod", "host", "evil.com", "user", "root", "aws"}, tags)
	})

	t.Run("prunes multiple tag comment lines on update", func(t *testing.T) {
		content := "Host test-host\n    # tags: tag1\n    HostName 127.0.0.1\n    # tags: tag2\n"
		cfg, err := ssh_config.Decode(strings.NewReader(content))
		assert.NoError(t, err)

		err = UpdateASTHostTags(cfg.Hosts[1], []string{"merged1", "merged2"})
		assert.NoError(t, err)

		tags := ExtractTagsFromNodes(cfg.Hosts[1].Nodes)
		assert.Equal(t, []string{"merged1", "merged2"}, tags)
	})
}
