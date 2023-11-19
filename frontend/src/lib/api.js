// const API_HOST = window.location.origin;
const API_HOST = import.meta.env.VITE_API_HOST;

const validate = async () => {
  let res = await fetch(`${API_HOST}/auth/validate`, {
    credentials: "include",
  });

  return res;
};

const login = async (body) => {
  let res = await fetch(`${API_HOST}/auth/login`, {
    method: "POST",
    body: JSON.stringify(body),
  });

  return res;
};

const logout = async () => {
  let res = await fetch(`${API_HOST}/auth/logout`, {
    credentials: "include",
  });

  return res;
};

const register = async (body) => {
  let res = await fetch(`${API_HOST}/auth/register`, {
    method: "POST",
    body: JSON.stringify(body),
  });

  return res;
};

const getItems = async () => {
  let res = await fetch(`${API_HOST}/api/items`, {
    credentials: "include",
  });

  return res;
};

const getItem = async (id) => {
  let res = await fetch(`${API_HOST}/api/items/${id}`, {
    credentials: "include",
  });

  return res;
};

const addItem = async (item) => {
  let res = await fetch(`${API_HOST}/api/items`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(item),
  });

  return res;
};

const deleteItem = async (item) => {
  let res = await fetch(`${API_HOST}/api/items`, {
    method: "DELETE",
    credentials: "include",
    body: JSON.stringify(item),
  });

  return res;
};

const addCollection = async (collection) => {
  let res = await fetch(`${API_HOST}/api/collections`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(collection),
  });

  return res;
};

const getCollections = async () => {
  let res = await fetch(`${API_HOST}/api/collections`, {
    credentials: "include",
  });

  return res;
};

const getCollection = async (id) => {
  let res = await fetch(`${API_HOST}/api/collections/${id}`, {
    credentials: "include",
  });

  return res;
};

const getCollectionItems = async (id) => {
  let res = await fetch(`${API_HOST}/api/items/${id}`, {
    credentials: "include",
  });

  return res;
};

const updateCollection = async (collection) => {
  let res = await fetch(`${API_HOST}/api/collections`, {
    method: "PUT",
    credentials: "include",
    body: JSON.stringify(collection),
  });

  return res;
};

const archiveCollection = async (collection) => {
  let res = await fetch(
    `${API_HOST}/api/collections/${collection.id}/archive`,
    {
      method: "POST",
      credentials: "include",
      body: JSON.stringify(collection),
    },
  );

  return res;
};

const deleteCollection = async (collection) => {
  let res = await fetch(`${API_HOST}/api/collections`, {
    credentials: "include",
    method: "DELETE",
    body: JSON.stringify(collection),
  });

  return res;
};

const getPageInfo = async (url) => {
  let res = await fetch(`${API_HOST}/api/opengraph/info`, {
    credentials: "include",
    method: "POST",
    body: JSON.stringify({
      url,
    }),
  });

  return res;
};

export default {
  validate,
  login,
  logout,
  register,

  getItems,
  getItem,
  addItem,
  deleteItem,

  getCollections,
  getCollection,
  getCollectionItems,
  addCollection,
  updateCollection,
  archiveCollection,
  deleteCollection,

  getPageInfo,
};
