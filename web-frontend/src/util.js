export function checkValidUrl(url) {
	if (!url || (!url.startsWith("http://") && !url.startsWith("https://"))) {
		console.log("invalid", url);
		return false;
	}
	return true;
}

export function displayTimestamp(unixMillis) {
	return new Date(unixMillis).toLocaleString();
}
