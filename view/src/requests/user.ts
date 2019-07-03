import XMLHttp, {HTTPMethods} from "../utils/xml_http";

export const create = (async (name) => {
  const json = JSON.stringify({name});
  return await XMLHttp.request("/create_user", {}, HTTPMethods.Post, json);
});

export const login = (async (loginToken) => {
  const header = { "x-authenticate-token": loginToken };
  return await XMLHttp.request("/login", header, HTTPMethods.Post, "");
});
