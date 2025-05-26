import { like, unlike, updateLikeCount, likeStatus } from "/api/like";

const urlParams = new URLSearchParams(window.location.search);
const postId = urlParams.get('id');

if (postId) {
  fetch(`/api/post/${postId}`)
    .then(response => response.json())
    .then(post => {
      renderPostDetail(post);
    })
    .catch(error => {
      console.error("Erreur lors du chargement du post :", error);
    });
} else {
  console.error("Aucun ID de post fourni.");
}


function renderPostDetail(post) {
  const postContainer = document.getElementById("dynamic-post");
  postContainer.innerHTML = `
  <div class="title-author">
    <h1>${post.Title}</h1>
    <p class="author">${post.Author || "Auteur inconnu"}</p>
  </div>
  <p class="content">${post.Content}</p>
  <div class="category">${post.Themes}</div>
  <div class="date">Publié le ${new Date(post.CreatedAt).toLocaleDateString("fr-FR")}</div>
  <div class="likes">Likes : ${post.Likes}</div>
  `;

  const likeBtn = document.querySelector('.like-btn');
  let isLiked = false;
  isLiked = likeStatus(post.ID);
  if (isLiked)
    likeBtn.classList.add('liked');

  likeBtn.addEventListener('click', async () => {
    if (isLiked) {
      await unlike(post.ID);
      isLiked = false;
      likeBtn.classList.remove('liked');
    } else {
      await like(post.ID);
      isLiked = true;
      likeBtn.classList.add('liked');
    }
    updateLikeCount(post.ID);
  });
}

document.addEventListener("DOMContentLoaded", async () => {
  const postContainer = document.getElementById("dynamic-post");
  const params = new URLSearchParams(window.location.search);
  const postId = params.get("id");

  if (!postContainer) {
    console.error("Élément #dynamic-post introuvable dans le DOM");
    return;
  }

  if (!postId) {
    postContainer.innerHTML = "<p>Post introuvable (ID manquant).</p>";
    return;
  }

  try {
    const res = await fetch(`/api/post/${postId}`);
    if (!res.ok) throw new Error("Post non trouvé");

    const post = await res.json();

    postContainer.innerHTML = "";
    postContainer.innerHTML = `
      <div class="title-author">
        <h1>${post.Title}</h1>
        <p class="author">${post.Author || "Auteur inconnu"}</p>
      </div>
      <p class="content">
        ${post.Content}
      </p>
      <div class="category">Catégorie : ${post.Themes}</div>
      <div class="date">Publié le ${new Date(post.CreatedAt).toLocaleDateString("fr-FR")}</div>
    `;
    console.log("Post chargé avec succès :", post);
    updateLikeCount(post.ID);

    const likeBtn = document.querySelector('.like-btn');
    let likeSpan = document.createElement('span');
    likeSpan.id = `like-count-${post.ID}`;
    likeSpan.textContent = post.Likes || 0;
    likeBtn.appendChild(likeSpan);

  } catch (err) {
    console.error("Erreur:", err);
    postContainer.innerHTML = "<p>Erreur lors du chargement du post.</p>";
  }
});

document.addEventListener("DOMContentLoaded", () => {
  const authContainer = document.querySelector(".auth-buttons");

  const buttons = [
    { text: "Connexion", url: "/front/login/login.html", id: "connexion-button" },
    { text: "Inscription", url: "/front/register/register.html", id: "inscription-button" },
    { text: "Profil", url: "/front/profil/profil.html", id: "profil-button" },
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

  const sessionCookie = document.cookie.split("; ").find(row => row.startsWith("session_token="));

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