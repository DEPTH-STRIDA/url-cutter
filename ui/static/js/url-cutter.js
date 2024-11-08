document.addEventListener("DOMContentLoaded", function () {});

// Проверяет валидность URL
const isValidUrl = (urlString) => {
  // Убираем пробелы в начале и конце строки
  urlString = urlString.trim();

  // Если URL начинается с www, добавляем протокол для корректной обработки
  if (urlString.startsWith("www.")) {
    urlString = "http://" + urlString;
  }

  // Проверяем, начинается ли строка с http:// или https://
  const urlPattern = /^(http:\/\/|https:\/\/)/;

  // Если URL не начинается с одной из разрешенных схем, проверяем на локальный адрес
  if (
    !urlPattern.test(urlString) &&
    !urlString.startsWith("localhost") &&
    !urlString.startsWith("127.0.0.1")
  ) {
    return false;
  }

  try {
    // Парсим URL и проверяем hostname
    const parsedUrl = new URL(urlString);
    return parsedUrl.hostname !== ""; // Проверяем, что hostname не пустой
  } catch (err) {
    return false; // Если произошла ошибка, возвращаем false
  }
};

// Проверка URL на валидность
document.addEventListener("DOMContentLoaded", function () {
  let inputField = document.getElementById("url-input");

  inputField.addEventListener("input", (event) => {
    const currentValue = event.target.value;

    if (isValidUrl(currentValue)) {
      enableButton();
    } else {
      disableMainButton();
    }
  });
});

// Активирует кнопку для отправки url
function enableButton() {
  let button = document.getElementById("url-button");
  button.classList.remove("disabled");
  button.disabled = false;
}
// Выключает кнопку для отправки url
function disableMainButton() {
  let button = document.getElementById("url-button");
  button.classList.add("disabled");
  button.disabled = true;
}

// Покас/скрытие контейнера с поделиться
let shareContainerStatus = false;

function toggleShareContainer() {
  const container = document.getElementById("share-container");
  if (shareContainerStatus) {
    container.style.height = "0";
    container.style.opacity = "0";
    setTimeout(() => {
      container.classList.remove("active");
    }, 300); // Время должно соответствовать длительности transition в CSS
    shareContainerStatus = false;
  } else {
    container.classList.add("active");
    const height = container.scrollHeight + "px";
    container.style.height = height;
    container.style.opacity = "1";
    shareContainerStatus = true;
  }
}

document.addEventListener("click", function (event) {
  const clickedElement = event.target;
  if (
    clickedElement.id == "share-button" ||
    clickedElement.id == "share-button-img"
  ) {
    toggleShareContainer();
    event.stopPropagation();
  } else if (
    shareContainerStatus &&
    !event.target.closest("#share-container")
  ) {
    toggleShareContainer();
  }
});

// Показ/скрытие значка загрузки
function enabledLoading() {
  let button = document.getElementById("url-button");
  button.style.color = "#222222";
  button.classList.add("loading");
}

function disabledLoading() {
  let button = document.getElementById("url-button");
  button.style.color = "white";
  button.classList.remove("loading");
}

// Показ/скрытие короткой ссылки
let isShortUrlVisible = false;

function enabledshortUrlContainer() {
  let shortUrlContainer = document.querySelector(".short-url");
  if (!isShortUrlVisible) {
    shortUrlContainer.classList.add("visible");
    isShortUrlVisible = true;
  }
}

function disabledshortUrlContainer() {
  let shortUrlContainer = document.querySelector(".short-url");
  if (isShortUrlVisible) {
    shortUrlContainer.classList.remove("visible");
    isShortUrlVisible = false;
  }
}

document.addEventListener("DOMContentLoaded", function () {
  const urlButton = document.getElementById("url-button");
  const urlInput = document.getElementById("url-input");

  urlButton.addEventListener("click", function () {
    const url = urlInput.value.trim();
    if (!isValidUrl(url)) {
      alert("Пожалуйста, введите корректный URL");
      return;
    }

    const data = { "long-url": url };
    const dataJSON = JSON.stringify(data);

    enabledLoading();
    disabledshortUrlContainer();

    sendPostRequestAsync("/shorten", dataJSON)
      .then((result) => {
        if (result.success) {
          const response = JSON.parse(result.response);
          if (response["short-url"] && response.img) {
            // Проверяем, есть ли уже контент в контейнере
            const urlContainer = document.querySelector(".url-container");
            const hasExistingContent = urlContainer.textContent.trim() !== "";

            if (hasExistingContent) {
              // Если контент уже есть, применяем задержку
              setTimeout(() => {
                updateShortUrlContainer(response["short-url"], response.img);
                enabledshortUrlContainer();
              }, 500);
            } else {
              // Если контента нет, обновляем и показываем сразу
              updateShortUrlContainer(response["short-url"], response.img);
              enabledshortUrlContainer();
            }
          } else {
            alert("Произошла ошибка при сокращении ссылки");
          }
        } else {
          alert("Ошибка сервера: " + result.response);
        }
      })
      .catch((error) => {
        alert("Произошла ошибка: " + error.message);
      })
      .finally(() => {
        disabledLoading();
      });
  });
});

function updateShortUrlContainer(shortUrl, imgBase64) {
  const urlContainer = document.querySelector(".url-container");
  const qrImage = document.querySelector(".qr img");

  urlContainer.textContent = shortUrl;
  qrImage.src = "data:image/png;base64," + imgBase64;
}

document.addEventListener("click", function (event) {
  const clickedElement = event.target;
  if (clickedElement.id == "copy-button") {
    let textContainer = document.getElementById("ready-short-url");
    let textToCopy = textContainer.textContent;

    navigator.clipboard
      .writeText(textToCopy)
      .then(() => {
        // alert("Текст скопирован в буфер обмена!");
      })
      .catch((err) => {
        console.error("Ошибка копирования: ", err);
      });
  }
});

// Скачивание картинки
document.addEventListener("click", function (event) {
  const clickedElement = event.target;
  if (clickedElement.id == "download-qr-code-button") {
    let img = document.getElementById("qr-code-img");
    let imgSrc = img.src;

    let textContainer = document.getElementById("ready-short-url");
    let textToCopy = textContainer.textContent;

    downloadImage(imgSrc, textToCopy);
  }
});

function downloadImage(base64Image, name) {
  // Создаём временный элемент <a>
  const link = document.createElement("a");
  link.href = base64Image;
  link.download = name + ".png"; // Имя файла для скачивания

  // Добавляем элемент в документ
  document.body.appendChild(link);
  // Инициируем клик на элементе <a>
  link.click();
  // Удаляем элемент после скачивания
  document.body.removeChild(link);
}

// Реализация кнопок "Поделиться"
document.addEventListener("click", function (event) {
  const clickedElement = event.target;

  // Получаем текст только один раз
  const textContainer = document.getElementById("ready-short-url");
  const textToCopy = textContainer.textContent;
  const encodedUrl = encodeURIComponent("https://" + textToCopy); // Кодируем URL для передачи
  const message = encodeURIComponent("Check this out!"); // Текст сообщения для ВКонтакте

  switch (clickedElement.id) {
    case "tg-share-button":
      let telegramUrl = `https://t.me/share/url?url=${encodedUrl}`;
      window.open(telegramUrl, "_blank");
      break;
    case "wa-share-button":
      let whatsappUrl = `https://wa.me/?text=${encodedUrl}`;
      window.open(whatsappUrl, "_blank");
      break;
    case "vb-share-button":
      let viberUrl = `viber://forward?text=${encodedUrl}`;
      window.open(viberUrl, "_blank");
      break;
    case "vk-share-button":
      // Используем правильно закодированный URL
      let vkUrl = `https://vk.com/share.php?url=${encodedUrl}&title=${message}`;
      window.open(vkUrl, "_blank");
      break;
  }
});
