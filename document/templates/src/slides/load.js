// TODO: implement compression
const decompress = (str) => {
  const ds = new DecompressionStream("gzip")
  const blob = new Blob([str], { type : 'plain/text' })
  const decompressedStream = blob.stream().pipeThrough(ds)
  return new Response(decompressedStream).blob() // Promise
}

const loadSlides = (elementID) => new Promise((resolve, reject) => {
  try {
    const element = document.getElementById(elementID)
    if (!element) throw new Error(`element "${elementID}" could not be found on the page`)

    const content = element.textContent
    if (!content) throw new Error(`element "${elementID}" text content is empty`)

    const slideData = JSON.parse(content)
    const result = {
      slides: [],
      notes: [],
      footnotes: []
    }
    if (!Array.isArray(slideData)) throw new Error(`element "${elementID}" text content does not contain a JSON array`)
    for (let i = 0; i < slideData.length; i += 3) {
      result.slides.push(slideData[i] || "")
      result.notes.push(slideData[i+1] || "")
      result.footnotes.push(slideData[i+2] || "")
    }

    if (import.meta.env.DEV) {
      for (let i = 0; i < 3; i++) {
        result.slides = result.slides.concat(result.slides)
        result.notes = result.notes.concat(result.notes)
        result.footnotes = result.footnotes.concat(result.footnotes)
      }
    }

    // resolve(result)
    setTimeout(() => resolve(result), 1000)
  } catch (e) {
    reject(e)
  }
})

export default loadSlides
