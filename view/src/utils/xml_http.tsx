export enum HTTPMethods {
  Get = "GET",
  Post = "POST",
  Delete = "DELETE",
  Put = "PUT",
}

interface IResponse {
  status: number;
  data: any;
}

export default class XMLHttp {
  public static request(url: string, headers: {}, method: HTTPMethods, data: string): Promise<IResponse> {
    return new Promise(((resolve, reject) => {
      const req = new XMLHttpRequest();
      const requestURL = `http://localhost:19999${url}`;
      req.open(method, requestURL, true);
      for (const key in headers) {
        if (headers.hasOwnProperty(key)) {
          req.setRequestHeader(key, headers[key]);
        }
      }
      req.onreadystatechange = () => {
        if (req.readyState === XMLHttpRequest.DONE) {
          resolve({ status: req.status, data: JSON.parse(req.responseText) });
        }
      };
      req.onerror = () => {
        reject(new Error(req.statusText));
      };
      req.send(data === "" ? null : data);
    }));
  }
}
