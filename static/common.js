(() => {
  console.log("TTS");
  const itemKey = "initSsml"
  const initValue = localStorage.getItem(itemKey)
  if (initValue) {
    document.getElementById("ssmlText").value = initValue
  }
  document.getElementById("synthesizeButton").addEventListener("click", () => {
    const value = document.getElementById("ssmlText").value
    localStorage.setItem(itemKey, value)
    const form = new FormData(document.getElementById("ssml"));
    const headers = {
      Accept: "application/json",
    };
    fetch("/api/synthesize", { method: "POST", headers, body: form })
      .then((res) => {
        let audio = document.getElementById("audio");
        audio.src = "/synthesize.mp3?" + Date.now();
        audio.play();
      })
      .catch(console.error);
  });
})();
