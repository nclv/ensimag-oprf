if (!WebAssembly.instantiateStreaming) { // polyfill
	WebAssembly.instantiateStreaming = async (resp, importObject) => {
		const source = await (await resp).arrayBuffer();
		return await WebAssembly.instantiate(source, importObject);
	};
}

const go = new Go();
let mod, inst;
WebAssembly.instantiateStreaming(fetch("/static/client.wasm"), go.importObject).then((result) => {
	mod = result.module;
	inst = result.instance;
	// activate the button on wasm file load
	document.getElementById("runButton").disabled = false;
}).then(() => {
	// load the pseudonymize function
	void run();
}).catch((err) => {
	console.error(err);
});

// run the WASM client code
async function run() {
	console.clear();
	console.log("Running the WASM");

	await go.run(inst);
	// instance reset can be commented
	inst = await WebAssembly.instantiate(mod, go.importObject); // reset instance
}

// pseudonymize the JSON input
async function AsyncPseudonymize(jsonInput) {
	try {
		const objectOutput = await pseudonymize(jsonInput);
		console.log(objectOutput);

		return objectOutput;
	} catch (err) {
		console.error('Caught exception', err);
	}
}

async function handleFormSubmit(event) {
	console.log("Form submitted");
	event.preventDefault();

	const formData = new FormData(event.target);
	const formObject = Object.fromEntries(formData.entries());

	console.log(formObject);

	formObject["mode"] = parseInt(formObject["mode"]);
	formObject["suite"] = formObject["suite"];
	formObject["return-info"] = formObject["return-info"] === 'true';
	formObject["data"] = formObject["data"].split(";");

	const objectOutputs = await AsyncPseudonymize(JSON.stringify(formObject));

	const results = document.querySelector('.results pre');
	results.innerText = JSON.stringify(objectOutputs, null, 4);
}

const load = () => {
	const form = document.querySelector('.pseudonymization-form');
	form.addEventListener('submit', (event) => handleFormSubmit(event));
}
window.onload = load;