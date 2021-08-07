const HelloHandler = (req, res) => {
	res.statusCode = 200

	res.json({
		text: 'Hello World! This is the Next.js starter kit :D',
	})
}

export default HelloHandler
