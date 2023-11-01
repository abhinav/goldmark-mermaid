async function renderSVG(src) {
	const { svg } = await mermaid.render('mermaid', src);
	return svg;
}
