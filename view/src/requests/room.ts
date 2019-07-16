import XMLHttp, {HTTPMethods} from "../utils/xml_http";

export const join = (async (roomId, loginToken, accessToken) => {
  const header = { "x-authenticate-token": loginToken, "x-access-token": accessToken };
  const json = JSON.stringify({room_id: roomId});
  return await XMLHttp.request("/join_room", header, HTTPMethods.Post, json);
});

export const message = (async (loginToken, accessToken, text) => {
  const header = { "x-authenticate-token": loginToken, "x-access-token": accessToken };
  const json = JSON.stringify({ text });
  return await XMLHttp.request("/message", header, HTTPMethods.Post, json);
});
