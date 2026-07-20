import {defineConfig} from 'oxfmt';

export default defineConfig({
	printWidth: 80,
	trailingComma: 'all',
	singleQuote: true,
	arrowParens: 'avoid',
	bracketSpacing: false,
	bracketSameLine: true,
	useTabs: true,
	sortTailwindcss: {
		stylesheet: './src/index.css',
		functions: ['clsx', 'cn'],
		preserveWhitespace: true,
	},
});
