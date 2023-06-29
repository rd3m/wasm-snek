/** @type {import('tailwindcss').Config} */
module.exports = {
	content: ["./cmd/server/templates/*.html"],
	theme: {
		extend: {},
	},
	daisyui: {
		themes: ["retro"],
	},

	plugins: [require("daisyui")],
};
