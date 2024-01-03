package proc

var rssTemplate = `<?xml version="1.0" encoding="utf-8"?>
<rss xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd" version="2.0">
  <channel>
	<title>{{.FeedTitle}}</title>
	<link>{{.FeedURL}}</link>
	<language>ru</language>
	<copyright>Umputun, Bobuk, Gray, Ksenks, Alek.sys</copyright>
	<itunes:author>Umputun, Bobuk, Gray, Ksenks, Alek.sys</itunes:author>
	<itunes:subtitle>{{.FeedSubtitle}}</itunes:subtitle>
	<description>{{.FeedDescription}}</description>
	<itunes:explicit>no</itunes:explicit>

	<itunes:owner>
		  <itunes:name>Umputun, Bobuk, Gray, Ksenks, Alek.sys</itunes:name>
		  <itunes:email>podcast@radio-t.com</itunes:email>
	</itunes:owner>

	<itunes:image href="{{.FeedImage}}" />

	<itunes:category text="Technology">
		  <itunes:category text="Tech News"/>
	</itunes:category>

	<itunes:keywords>hitech,russian,radiot,tech,news,радио</itunes:keywords>

	{{- range .Items}}
	<item>
		  <title>{{.Title}}</title>
          {{- if .Description}}
		  <description><![CDATA[{{- .Description -}}]]></description>
		  {{- end}}
		  <link>{{.URL}}</link>
		  <guid>{{.GUID}}</guid>
		  <pubDate>{{.Date}}</pubDate>
		  <itunes:author>Umputun, Bobuk, Gray, Ksenks, Alek.sys</itunes:author>
		  {{- if .Summary}}
		  <itunes:summary><![CDATA[{{.Summary -}}]]></itunes:summary>
		  {{- end}}
		  <itunes:image href="{{.Image}}" />
		  <enclosure url="{{.EnclosureURL}}" type="audio/mp3" {{ if .FileSize -}} length="{{.FileSize}}"{{- end }} />
	</item>
	{{- end}}
  </channel>
</rss>
`
