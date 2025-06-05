document.addEventListener("DOMContentLoaded", async () => {
  const authContainer = document.querySelector(".auth-buttons");
  const sessionCookie = document.cookie.split("; ").find(row => row.startsWith("session_token="));

  let username
  if (sessionCookie) {
    await fetch("/get-profile", {
      method: "GET",
      credentials: "include",
    })
      .then(response => {
        if (response.ok) {
          return response.json();
        } else {
          throw new Error("Erreur lors de la récupération du profil");
        }
      })
      .then(data => {
        username = data.profile.username;
      })
      .catch(error => {
        console.error("Erreur:", error);
      });
  }

  const buttons = [
    { text: "Connexion", url: "/front/login/login.html", id: "connexion-button" },
    { text: "Inscription", url: "/front/register/register.html", id: "inscription-button" },
    { text: "Profil", url: `/front/profil/profil.html?user=${username}`, id: "profil-button" },
  ];

  buttons.forEach(btn => {
    const button = document.createElement("button");
    button.className = "auth-button";
    button.textContent = btn.text;
    button.id = btn.id;
    button.onclick = () => {
      window.location.href = btn.url;
    };
    authContainer.appendChild(button);
  });


  if (sessionCookie) {
    document.getElementById("connexion-button").style.display = "none";
    document.getElementById("inscription-button").style.display = "none";
  } else {
    document.getElementById("profil-button").style.display = "none";
  }

  const ctaButton = document.querySelector(".cta-button");
  if (ctaButton) {
    ctaButton.addEventListener("click", () => {
      window.location.href = "/front/post-list/postlist.html";
    });
  }
});
