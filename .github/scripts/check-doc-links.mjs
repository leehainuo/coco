import { execFileSync } from 'node:child_process'
import { existsSync, readFileSync, statSync } from 'node:fs'
import { dirname, extname, isAbsolute, relative, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

// Ignore code examples while preserving line numbers in diagnostics.
function prose(markdown) {
  let fence
  return markdown.split('\n').map((line) => {
    const marker = line.match(/^\s{0,3}(`{3,}|~{3,})/)
    if (marker && !fence) {
      fence = marker[1]
      return ''
    }
    if (fence) {
      if (marker && marker[1][0] === fence[0] && marker[1].length >= fence.length) fence = undefined
      return ''
    }
    return line
  }).join('\n')
}

function anchors(markdown) {
  const result = new Set()
  const counts = new Map()
  for (const match of prose(markdown).matchAll(/^#{1,6}\s+(.+?)\s*#*\s*$/gm)) {
    const slug = match[1].toLowerCase()
      .replace(/<[^>]*>/g, '')
      .replace(/!?\[([^\]]+)\]\([^)]*\)/g, '$1')
      .replace(/[^\p{L}\p{M}\p{N}_\-\s]/gu, '')
      .replace(/\s/g, '-')
    let anchor = slug
    let count = counts.get(slug) ?? 0
    while (result.has(anchor)) anchor = `${slug}-${++count}`
    counts.set(slug, count)
    result.add(anchor)
  }
  for (const match of markdown.matchAll(/\b(?:id|name)=["']([^"']+)["']/g)) result.add(match[1])
  return result
}

export function checkLinks(root, files) {
  const errors = []
  const anchorCache = new Map()
  for (const file of files) {
    const source = resolve(root, file)
    const text = prose(readFileSync(source, 'utf8'))
    const links = [
      ...text.matchAll(/!?\[[^\n]*?\]\(\s*(<[^>]+>|[^\s)]+)(?:\s+["'][^\n]*?["'])?\s*\)/g),
      ...text.matchAll(/\b(?:href|src)=["']([^"']+)["']/g),
      ...text.matchAll(/^\s{0,3}\[[^\]]+\]:\s*(<[^>]+>|\S+)/gm),
    ]
    for (const match of links) {
      const target = match[1].replace(/^<|>$/g, '')
      if (/^(?:[a-z][a-z\d+.-]*:|\/\/)/i.test(target)) continue
      const line = text.slice(0, match.index).split('\n').length
      let path, fragment
      try {
        const [urlPath, ...hash] = target.split('#')
        path = decodeURIComponent(urlPath.split('?')[0])
        fragment = decodeURIComponent(hash.join('#'))
      } catch {
        errors.push(`${file}:${line}: invalid URL encoding: ${target}`)
        continue
      }
      const destination = path ? resolve(dirname(source), path) : source
      const local = relative(root, destination)
      if (local === '..' || local.startsWith('../') || isAbsolute(local) || !existsSync(destination)) {
        errors.push(`${file}:${line}: missing repository path: ${target}`)
        continue
      }
      if (fragment && extname(destination).toLowerCase() === '.md' && statSync(destination).isFile()) {
        if (!anchorCache.has(destination)) anchorCache.set(destination, anchors(readFileSync(destination, 'utf8')))
        if (!anchorCache.get(destination).has(fragment)) errors.push(`${file}:${line}: missing heading: ${target}`)
      }
    }
  }
  return errors
}

if (process.argv[1] && resolve(process.argv[1]) === fileURLToPath(import.meta.url)) {
  const root = resolve(dirname(fileURLToPath(import.meta.url)), '../..')
  const files = execFileSync('git', ['ls-files', '-z', '*.md'], { cwd: root, encoding: 'utf8' }).split('\0').filter(Boolean)
  const errors = checkLinks(root, files)
  if (errors.length) {
    console.error(errors.join('\n'))
    process.exitCode = 1
  } else {
    console.log(`Local documentation links and heading anchors passed (${files.length} Markdown files).`)
  }
}
