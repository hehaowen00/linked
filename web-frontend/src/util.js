export function checkValidUrl(url) {
	let testUrl;
	try {
		testUrl = new URL(url);
		return true;
	} catch (e) {
		console.log("invalid url", url);
		return false;
	}
}

export function displayTimestamp(unixMillis) {
	return new Date(unixMillis).toLocaleString();
}
