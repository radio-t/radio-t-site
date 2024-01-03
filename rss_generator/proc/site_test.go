package proc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHugoSite_ReadPosts(t *testing.T) {

	t.Run("valid location", func(t *testing.T) {
		s := HugoSite{Location: "testdata/hugo", BaseURL: "https://radio-t.com"}
		posts, err := s.ReadPosts()
		require.NoError(t, err)
		require.Equal(t, 3, len(posts))

		assert.Equal(t, "2023-12-23 16:47:22 +0000 UTC", posts[0].CreatedAt.String())
		assert.Equal(t, "https://radio-t.com/p/2023/12/23/podcast-889/", posts[0].URL)
		assert.Contains(t, posts[0].Data, "[Apple останавливает продажи часов]")
		assert.Equal(t, "Радио-Т 889", posts[0].Config["title"].(string))
		assert.Equal(t, "https://radio-t.com/images/radio-t/rt889.jpg", posts[0].Config["image"].(string))
		assert.Equal(t, "rt_podcast889", posts[0].Config["filename"].(string))
		assert.Equal(t, []any{"podcast"}, posts[0].Config["categories"])

		assert.Equal(t, "2023-08-01 14:45:18 +0000 UTC", posts[1].CreatedAt.String())
		assert.Equal(t, "https://radio-t.com/p/2023/08/01/prep-870/", posts[1].URL)
		assert.Equal(t, "Темы для 870", posts[1].Config["title"].(string))

		assert.Equal(t, "2023-03-18 18:24:29 +0000 UTC", posts[2].CreatedAt.String())
		assert.Equal(t, "https://radio-t.com/p/2023/03/18/podcast-850/", posts[2].URL)
		assert.Equal(t, "Радио-Т 850", posts[2].Config["title"].(string))
	})

	t.Run("invalid location", func(t *testing.T) {
		s := HugoSite{Location: "testdata/invalid", BaseURL: "https://radio-t.com"}
		_, err := s.ReadPosts()
		require.Error(t, err)
	})
}
