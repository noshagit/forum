import { updateLikeCount, like, unlike } from '/api/like';

const postsContainer = document.getElementById("posts-container");

async function fetchPosts() {
  try {
    const res = await fetch("/api/posts");
    const posts = await res.json();

    console.log("Résultat brut du fetch:", posts);

    if (!Array.isArray(posts) || posts.length === 0) {
      console.warn("Aucun post trouvé ou format inattendu:", posts);
    } else {
      console.log(`Nombre de posts récupérés: ${posts.length}`);
    }

    posts.forEach(post => {
      console.log("Post individuel:", post);

      const wrapper = document.createElement("div");
      wrapper.className = "post-container";

      const postEl = document.createElement("div");
      postEl.className = "post";

      postEl.innerHTML = `
    <h2>${post.Title}</h2>
    <p>${post.Content}</p>
    <div class="category">${post.Themes || "Aucune catégorie"}</div>
    <button class="show-more">Voir plus</button>
  `;

      const reactions = document.createElement("div");
      reactions.className = "reactions";

      reactions.innerHTML = `
    <div class="reaction-box like-btn" id="${post.ID}">
      <img src="/front/images/like.png" alt="like">
      <span class="like-count" id="like-count-${post.ID}">${post.Likes}</span>
    </div>
    <div class="reaction-box"><img src="/front/images/share.png" alt="share"></div>
  `;

      wrapper.appendChild(postEl);
      wrapper.appendChild(reactions);
      postsContainer.appendChild(wrapper);

      const showMoreBtn = postEl.querySelector(".show-more");
      showMoreBtn.addEventListener("click", () => {
        postEl.classList.toggle("expanded");
        showMoreBtn.textContent = postEl.classList.contains("expanded") ? "Voir moins" : "Voir plus";
      });

      const likeBtn = reactions.querySelector('.like-btn');
      let isLiked = false;

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

      updateLikeCount(post.ID);
    });


  } catch (err) {
    console.error("Erreur lors du chargement des posts:", err);
  }
}

document.addEventListener("DOMContentLoaded", () => {
  fetchPosts();
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
