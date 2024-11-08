function sendPostRequestAsync(url, data) {
  return new Promise((resolve, reject) => {
    var Request = false;

    if (window.XMLHttpRequest) {
      Request = new XMLHttpRequest();
    } else if (window.ActiveXObject) {
      try {
        Request = new ActiveXObject("Microsoft.XMLHTTP");
      } catch (CatchException) {
        try {
          Request = new ActiveXObject("Msxml2.XMLHTTP");
        } catch (CatchException2) {
          Request = false;
        }
      }
    }

    if (!Request) {
      reject(new Error("Невозможно создать XMLHttpRequest"));
      return;
    }

    console.log("URL:", url);
    console.log("Data:", data);

    Request.open("POST", url, true);

    Request.setRequestHeader(
      "Content-Type",
      "application/x-www-form-urlencoded"
    );

    Request.onreadystatechange = function () {
      if (Request.readyState === 4) {
        if (Request.status === 200) {
          resolve({ success: true, response: Request.responseText });
        } else {
          let response;
          switch (Request.status) {
            case 404:
              response = "404 (Not Found)";
              break;
            case 403:
              response = "403 (Forbidden)";
              break;
            case 500:
              response = "500 (Internal Server Error)";
              break;
            default:
              response = `${Request.status} (${
                Request.statusText || "Unknown Error"
              })`;
          }
          resolve({ success: false, response: response });
        }
      }
    };

    Request.onerror = function () {
      reject(new Error("Network Error"));
    };

    try {
      Request.send(data);
    } catch (error) {
      reject(error);
    }
  });
}
