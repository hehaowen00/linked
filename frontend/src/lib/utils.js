export function isValidURL(url) {
  try {
    let u = new URL(url);
    return true;
  } catch (e) {
    return false;
  }
}
