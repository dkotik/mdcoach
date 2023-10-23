import { fade } from 'svelte/transition'
import { elasticOut } from 'svelte/easing'
import currentSlide from './current.js'

let lastSlide = 0
let goingRight = true
currentSlide.subscribe(value => {
  goingRight = lastSlide < value
  lastSlide = value
})

const slideIn = (node, { duration }) => {
	return {
		duration,
		css: (t) => {
			const eased = elasticOut(t);
      return `transform: translate(${(goingRight ? 100 : - 100)*eased}%, 0%);`
      // console.log("goingRight", goingRight)
			// return `
			// 	transform: scale(${eased}) rotate(${eased * 1080}deg);
			// 	color: hsl(
			// 		${~~(t * 360)},
			// 		${Math.min(100, 1000 - 1000 * t)}%,
			// 		${Math.min(50, 500 - 500 * t)}%
			// 	);`;
		}
	};
}

const slideOut = (node, { duration }) => {
	return {
		duration,
		css: (t) => {
			const eased = elasticOut(t);
      return `transform: translate(${(goingRight ? -100 : 100)*eased}%, 0%); position: absolute;`
			// return `
			// 	transform: scale(${eased}) rotate(${eased * 1080}deg);
			// 	color: hsl(
			// 		${~~(t * 360)},
			// 		${Math.min(100, 1000 - 1000 * t)}%,
			// 		${Math.min(50, 500 - 500 * t)}%
			// 	);`;
		}
	};
}

export { slideIn, slideOut }
