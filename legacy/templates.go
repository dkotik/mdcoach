package mdcoach

// TODO: move this into assets?

const (
	tmplHead = `<!DOCTYPE html>
    <html xmlns="http://www.w3.org/1999/xhtml"><head>
      <meta http-equiv="content-type" content="text/html; charset=utf-8" />
      <title>{{ .title }}</title>
      <meta name="title" content="{{ .title }}" />
      <meta name="description" content="{{ .description }}" />
      <meta name="viewport" content="width=device-width, initial-scale=1.0" />
      <meta name="keywords" content="{{ .keywords }}" />
      <meta name="copyright" content="&copy; {{ .copyright }}" />
      <meta name="author" content="{{ .author }}" />
      <meta name="LastModifiedDate" content="{{ .updated }}" />
      <meta name="CreationDate" content="{{ .created }}" />
      <meta name="generator" content="MDCoach Beta" />
      {{ range .stylesheets }}
        <link rel="stylesheet" type="text/css" href="{{ . }}" />
      {{ end }}
      {{ range .scripts }}
        <script type="text/javascript" src="{{ . }}"></script>
      {{ end }}
    </head><body>
    <main>`
	// <link rel="stylesheet" type="text/css" href="pygments.css" />
	// <link rel="stylesheet" type="text/css" href="emote.css" />

	tmplPresentation = `<h4 style="text-align: center; color: #676767">&hellip;</h4>'];</script><main role="main">
      <div v-if="curtain" class="curtain"></div>
      <div v-else-if="navigator" class="navigator" v-cloak>
        <h1 v-text="document.title"></h1>
        <div class="notes">
          <span v-for="slide, key in slides">
            <a class="button" v-text="key+1" @click="window.location.hash='slide'+(key+1)"></a>
            &nbsp;
          </span>
        </div>
      </div>
      <a @click="slide-=1" class="leftMargin"></a>
      <a @click="slide+=1" class="rightMargin"></a>
      <section is="slide" v-for="slide, key in slides" :slideid="key"></section>
      <div class="controls" @mouseenter="controls=true" @mouseleave="controls=false">
        <nav v-if="controls">
          <!-- <a @click="alert('notes')"><svg viewBox="0 0 1792 1792" xmlns="http://www.w3.org/2000/svg"><path d="M448 1536h896v-256h-896v256zm0-640h896v-384h-160q-40 0-68-28t-28-68v-160h-640v640zm1152 64q0-26-19-45t-45-19-45 19-19 45 19 45 45 19 45-19 19-45zm128 0v416q0 13-9.5 22.5t-22.5 9.5h-224v160q0 40-28 68t-68 28h-960q-40 0-68-28t-28-68v-160h-224q-13 0-22.5-9.5t-9.5-22.5v-416q0-79 56.5-135.5t135.5-56.5h64v-544q0-40 28-68t68-28h672q40 0 88 20t76 48l152 152q28 28 48 76t20 88v256h64q79 0 135.5 56.5t56.5 135.5z"/></svg></a> -->
          <a @click="navigator=!navigator"><svg viewBox="0 0 1792 1792" xmlns="http://www.w3.org/2000/svg"><path d="M512 0q13 0 22.5 9.5t9.5 22.5v1472q0 20-17 28l-480 256q-7 4-15 4-13 0-22.5-9.5t-9.5-22.5v-1472q0-20 17-28l480-256q7-4 15-4zm1248 0q13 0 22.5 9.5t9.5 22.5v1472q0 20-17 28l-480 256q-7 4-15 4-13 0-22.5-9.5t-9.5-22.5v-1472q0-20 17-28l480-256q7-4 15-4zm-1120 0q8 0 14 3l512 256q18 10 18 29v1472q0 13-9.5 22.5t-22.5 9.5q-8 0-14-3l-512-256q-18-10-18-29v-1472q0-13 9.5-22.5t22.5-9.5z"/></svg></a>
          <a @click="toggleContrast"><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512"><path d="M8 256c0 136.966 111.033 248 248 248s248-111.034 248-248S392.966 8 256 8 8 119.033 8 256zm248 184V72c101.705 0 184 82.311 184 184 0 101.705-82.311 184-184 184z"/></svg></a>
          <a @click="curtain=!curtain"><svg viewBox="0 0 1792 1792" xmlns="http://www.w3.org/2000/svg"><path d="M1664 416v960q0 119-84.5 203.5t-203.5 84.5h-960q-119 0-203.5-84.5t-84.5-203.5v-960q0-119 84.5-203.5t203.5-84.5h960q119 0 203.5 84.5t84.5 203.5z"/></svg></a>
          <a @click="toggleFullScreen"><svg viewBox="0 0 1792 1792" xmlns="http://www.w3.org/2000/svg"><path d="M883 1056q0 13-10 23l-332 332 144 144q19 19 19 45t-19 45-45 19h-448q-26 0-45-19t-19-45v-448q0-26 19-45t45-19 45 19l144 144 332-332q10-10 23-10t23 10l114 114q10 10 10 23zm781-864v448q0 26-19 45t-45 19-45-19l-144-144-332 332q-10 10-23 10t-23-10l-114-114q-10-10-10-23t10-23l332-332-144-144q-19-19-19-45t19-45 45-19h448q26 0 45 19t19 45z"/></svg></a>
        </nav>
      </div>
    </main><script type="text/javascript" src=".cache/presentation.js"></script>
    </body></html>`

	tmplAssessment = `
    <table class="grades">
      <tr>
        <td colspan="8">&nbsp;</td>
        <th colspan="3">
          Группа:
        </th>
        <th colspan="6">
          Имя:
        </th>
      </tr>
      <tr>
        <td>&lt;</td>
        {{ range .missed }}
          <td>-{{ . }}&nbsp;</td>
        {{ end }}
      </tr>
      <tr>
        <td>0%</td>
        {{ range .grades }}
          <td>{{ . }}%</td>
        {{ end }}
      </tr>
    </table>
    <h1>{{ .title }}</h1>
    <section class="content">
    {{ .content }}
    <p>
      Ответьте на каждый вопрос кратко, точно, и чётко, чтобы получить +3 балла за каждый ответ.
    </p>
    </section>
    <section class="questions">
    <ol>
      {{ range $question := .questions }}
        <li>
          <legend>-1<br />-2<br />-3</legend>
          {{ $question }}
        </li>
      {{ end }}
    </ol>
    </section>
    </main><footer>
      {{ .footer }}
    </footer></body></html>`

	tmplReview = `{{ with .Questions }}<ol class="review">
      {{ range $question := . }}
        <li>{{ $question }}</li>
      {{ end }}
    </ol>{{ end }}`

	// tmplNotes = `<aside class="meta note content">
	//   <h1>{{ .title }}</h1>
	//   {{ if .description }}{{ .description }}<br />{{ end }}
	//   <small>[{{ .source }}{{ if .created }}, {{ .created }}{{ end }}]</small>
	// </aside>`

	// tmplMeta = `<aside class="meta note content">
	//   <h1>{{ .title }}</h1>
	//   {{ if .description }}{{ .description }}<br />{{ end }}
	//   <small>[{{ .source }}{{ if .created }}, {{ .created }}{{ end }}]</small>
	// </aside>`

	tmplIndex = `
    {{ with .syllabus }}
    {{ if .print }}<a target="_blank" href="{{ .print }}" class="button syllabus">&nbsp;</a>{{ end }}
    <h1>{{ .title }}</h1>
    <p>{{ .description }}</p>
    {{ end }}
    <ol>
      {{ range .sections }}
      <li>
        {{ if .view }}
        <a href="{{ .view }}" target="_blank" class="button presentation">&nbsp;</a>
        {{ end }}
        <a href="{{ .print }}" target="_blank" class="button notes">&nbsp;</a>
        <h2>{{ .title }}</h2>
        {{ .description }}
      </li>
      {{ end }}
    </ol>`

	tmplFoot = `</main></body></html>`
)
