<template>
  <div class="markdown-body" v-html="renderedContent"></div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'

const props = defineProps<{
  content: string
}>()

const md = new MarkdownIt({
  html: false,
  linkify: true,
  typographer: true,
  highlight: function (str: string, lang: string) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return `<pre class="hljs"><code>${hljs.highlight(str, { language: lang, ignoreIllegals: true }).value}</code></pre>`
      } catch (__) {}
    }
    return `<pre class="hljs"><code>${md.utils.escapeHtml(str)}</code></pre>`
  },
})

const renderedContent = computed(() => {
  if (!props.content) return ''
  return md.render(props.content)
})
</script>

<style>
.markdown-body {
  font-size: 15px;
  line-height: 1.7;
  color: #333;
  word-wrap: break-word;
}

.markdown-body p {
  margin: 1em 0;
}

.markdown-body h1,
.markdown-body h2,
.markdown-body h3,
.markdown-body h4,
.markdown-body h5,
.markdown-body h6 {
  margin-top: 1.5em;
  margin-bottom: 0.5em;
  font-weight: 600;
  line-height: 1.25;
  color: #1a1a1a;
}

.markdown-body h1 {
  font-size: 1.8em;
  padding-bottom: 0.3em;
  border-bottom: 1px solid #eee;
}

.markdown-body h2 {
  font-size: 1.5em;
  padding-bottom: 0.3em;
  border-bottom: 1px solid #eee;
}

.markdown-body h3 {
  font-size: 1.25em;
}

.markdown-body h4 {
  font-size: 1em;
}

.markdown-body ul,
.markdown-body ol {
  padding-left: 2em;
  margin: 1em 0;
}

.markdown-body li {
  margin: 0.25em 0;
}

.markdown-body li > p {
  margin: 0.5em 0;
}

.markdown-body blockquote {
  margin: 1em 0;
  padding: 0.5em 1em;
  border-left: 4px solid #ddd;
  background: #f9f9f9;
  color: #666;
}

.markdown-body blockquote p {
  margin: 0.5em 0;
}

.markdown-body code {
  padding: 0.2em 0.4em;
  margin: 0;
  font-size: 85%;
  background: #f1f1f1;
  border-radius: 4px;
  font-family: 'Fira Code', 'Courier New', monospace;
}

.markdown-body pre {
  margin: 1em 0;
  padding: 0;
  background: #1e1e1e;
  border-radius: 8px;
  overflow-x: auto;
}

.markdown-body pre code {
  display: block;
  padding: 1em;
  background: #1e1e1e;
  color: #d4d4d4;
  font-size: 14px;
  line-height: 1.6;
  overflow-x: auto;
}

.markdown-body table {
  width: 100%;
  margin: 1em 0;
  border-collapse: collapse;
}

.markdown-body table th,
.markdown-body table td {
  padding: 8px 12px;
  border: 1px solid #ddd;
}

.markdown-body table th {
  background: #f9f9f9;
  font-weight: 600;
}

.markdown-body table tr:nth-child(2n) {
  background: #f9f9f9;
}

.markdown-body hr {
  margin: 2em 0;
  border: none;
  border-top: 1px solid #ddd;
}

.markdown-body img {
  max-width: 100%;
  height: auto;
  border-radius: 4px;
}

.markdown-body a {
  color: #3b82f6;
  text-decoration: none;
}

.markdown-body a:hover {
  text-decoration: underline;
}

.markdown-body strong {
  font-weight: 600;
  color: #1a1a1a;
}

.markdown-body em {
  font-style: italic;
}

/* 代码高亮主题 */
.hljs {
  background: #1e1e1e;
  color: #d4d4d4;
}

.hljs-keyword,
.hljs-selector-tag,
.hljs-literal,
.hljs-section,
.hljs-link {
  color: #569cd6;
}

.hljs-string,
.hljs-title,
.hljs-name,
.hljs-type,
.hljs-attribute,
.hljs-symbol,
.hljs-bullet,
.hljs-addition,
.hljs-variable,
.hljs-template-tag,
.hljs-template-variable {
  color: #ce9178;
}

.hljs-comment,
.hljs-quote,
.hljs-deletion,
.hljs-meta {
  color: #6a9955;
}

.hljs-keyword,
.hljs-selector-tag,
.hljs-literal,
.hljs-title,
.hljs-section,
.hljs-doctag,
.hljs-type,
.hljs-name,
.hljs-strong {
  font-weight: bold;
}

.hljs-params,
.hljs-built_in {
  color: #9cdcfe;
}

.hljs-number {
  color: #b5cea8;
}

.hljs-function {
  color: #dcdcaa;
}
</style>
