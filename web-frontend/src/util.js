export function checkValidUrl(url) {
	let testUrl;
	try {
		testUrl = new URL(url);
		return true;
	} catch (e) {
		return false;
	}
}

export function displayTimestamp(unixMillis) {
	return new Date(unixMillis).toLocaleString();
}
