import { updateLikeCount, like, unlike, likeStatus } from '/api/like';

let allPosts = [];

async function fetchPosts() {
  try {
    const res = await fetch("/api/posts");
    const posts = await res.json();
    allPosts = posts;

    renderPosts(posts);

  } catch (err) {
    console.error("Erreur lors du chargement des posts:", err);
  }
}

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

function renderPosts(posts) {
  const postsContainer = document.getElementById("posts-container");
  postsContainer.innerHTML = "";

  posts.forEach(post => {
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
    `;

    wrapper.appendChild(postEl);
    wrapper.appendChild(reactions);
    postsContainer.appendChild(wrapper);

    const showMoreBtn = postEl.querySelector(".show-more");
    showMoreBtn.addEventListener("click", () => {
      window.location.href = `/front/comments/comments.html?id=${post.ID}`;
    });


    const likeBtn = reactions.querySelector('.like-btn');
    let isLiked = false;
    if (document.cookie.includes("session_token="))
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

    updateLikeCount(post.ID);
  });
}

document.addEventListener("DOMContentLoaded", () => {
  fetchPosts();

  const searchInput = document.querySelector(".search-bar");
  const categoryFilter = document.querySelector(".category-filter");
  const sortFilter = document.querySelector(".sort-filter");

  searchInput.addEventListener("input", applyFiltersAndSort);
  categoryFilter.addEventListener("change", applyFiltersAndSort);
  sortFilter.addEventListener("change", applyFiltersAndSort);
});

function applyFiltersAndSort() {
  const searchQuery = document.querySelector(".search-bar").value.toLowerCase();
  const selectedCategory = document.querySelector(".category-filter").value;
  const selectedSort = document.querySelector(".sort-filter").value;

  let filtered = allPosts.filter(post => {
    const matchSearch = post.Title.toLowerCase().includes(searchQuery);
    const matchCategory = selectedCategory === "" || post.Themes === selectedCategory;
    return matchSearch && matchCategory;
  });

  console.log(filtered.map(p => ({ id: p.id, likes: p.Likes || p.likes })));

  if (selectedSort === "date") {
    filtered.sort((a, b) => new Date(b.CreatedAt) - new Date(a.CreatedAt));
  } else if (selectedSort === "oldest") {
    filtered.sort((a, b) => new Date(a.CreatedAt) - new Date(b.CreatedAt));
  } else if (selectedSort === "likes") {
    filtered.sort((a, b) => b.Likes - a.Likes);
  } else if (selectedSort === "less-likes") {
    filtered.sort((a, b) => a.Likes - b.Likes);
  }

  renderPosts(filtered);
}

document.addEventListener("DOMContentLoaded", () => {
  const modal = document.getElementById("post-modal");
  const openBtn = document.getElementById("create-post-button");
  const closeBtn = modal.querySelector(".close-button");
  const form = document.getElementById("create-post-form");

  openBtn.addEventListener("click", () => {
    modal.classList.remove("hidden");
  });

  closeBtn.addEventListener("click", () => {
    modal.classList.add("hidden");
  });

  window.addEventListener("click", (e) => {
    if (e.target === modal) modal.classList.add("hidden");
  });

  form.addEventListener("submit", async (e) => {
    e.preventDefault();
    const formData = new FormData(form);

    try {
      const postData = {
        title: formData.get("title"),
        content: formData.get("content"),
        themes: formData.get("themes")
      };

      console.log(postData);

      const response = await fetch("/api/add-post", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(postData),
        credentials: "include",
      });


      if (response.ok) {
        alert("Post créé !");
        form.reset();
        modal.classList.add("hidden");

        await fetchPosts();
        applyFiltersAndSort();

      } else {
        alert("Erreur lors de la création du post.");
      }

    } catch (err) {
      console.error("Erreur:", err);
    }
  });
});

