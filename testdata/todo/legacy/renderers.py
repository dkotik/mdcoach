seoTrash = set([
    'i', 'a', 'about', 'an', 'are', 'as', 'at', 'be', 'by'', in', 'is', 'it', 'of', 'on',
    'or', 'that', 'the', 'this', 'to', 'was', 'what', 'when', 'where', 'who', 'will', 'with', 'the'])

def slugify(text, default='void'):
    return rawSlug(text, max_length=64, word_boundary=True, save_order=True, stopwords=seoTrash) or default

def summarizeUrl(url, length=24, cutoff=32):
    if len(url) <= length:
        return url
    pos = url.find(u'/', 10, cutoff)
    return url[:cutoff - 2 if pos == -1 else pos] + '/..'

class CustomBlockGrammar(BlockGrammar):
    thumbsup = re.compile(r'^(ok|kk|yes|ya|y|k|\+1)(?:[\s\Z]) *\n*')
    thumbsdown = re.compile(ur'^(no|nope|nah|neh|n|нет|н|\-1)(?:[\s\Z]) *\n*')
    bump = re.compile(ur'^(bump|бамп)(?:[\s\Z]) *\n*')

    def image(self, src, title, alt_text):
        if src.startswith('iframe '):
            m = re.search(r'src\=\"([^\"]+)\"', src)
            if m:
                src = m.group(1)
                if src.startswith('https://www.google.com/maps/embed?pb='):
                    return self.iframe(src, title, alt_text)
        elif src.startswith('http://2gis.') or src.startswith('http://go.2gis.com/km1h0'):
            return self.iframe(src, title, alt_text)
        elif src.startswith('https://www.youtube.com/watch?v='):
            return self.iframe(u'https://www.youtube.com/embed/' + src[32:], title, alt_text)
        if alt_text:
            return u'<figure>{}<figcaption>{}</figcaption></figure>'.format(
                super(Renderer, self).image(src, title, alt_text), escape(alt_text))
        else:
            return u'<figure class="nocaption">{}</figure>'.format(
                super(Renderer, self).image(src, title, alt_text))

    def block_emote(self, emoteClass, text):
        return u'<p class="grandemote grandemote-{}">{}</p>'.format(emoteClass, text)

    def block_citation(self, text):
        return u'<div class="citation">{}</div>'.format(text)
