import { like, unlike, updateLikeCount, likeStatus } from "/api/like";

const urlParams = new URLSearchParams(window.location.search);
const postId = urlParams.get('id');

const user = {
  id: null,
  name: null,
  email: null,
  avatar: null,
  isLoggedIn: false,

  async getProfile() {
    if (!document.cookie.includes("session_token="))
      return;
    try {
      const resp = await fetch("/get-profile", {
        method: "GET",
        credentials: "include"
      })
      if (!resp.ok) {
        console.error("Erreur lors de la récupération du profil :", resp.statusText);
        return;
      }
      const data = await resp.json();
      this.id = data.profile.id;
      this.name = data.profile.username;
      this.email = data.profile.email;
      this.avatar = data.profile.profile_picture;
      this.isLoggedIn = true;
    } catch (error) {
      console.error("Erreur lors de la récupération du profil :", error);
      this.isLoggedIn = false;
    }
  }
};

function setDocumentTitle(title) {
  if (title) {
    document.title = title;
  }
}

const shareBtn = document.querySelector('.reaction-box img[alt="Share"]');
if (shareBtn) {
  shareBtn.onclick = () => {
    const postUrl = window.location.href;
    navigator.clipboard.writeText(postUrl)
      .then(() => {
        alert("Lien copié dans le presse-papiers !");
      })
      .catch(() => {
        alert("Impossible de copier le lien.");
      });
  };
}

async function renderPostDetail(post) {
  const postContainer = document.getElementById("dynamic-post");
  postContainer.innerHTML = `
  <div class="title-author">
    <h1>${post.Title}</h1>
    <p class="author">${post.Author || "Auteur inconnu"}</p>
  </div>
  <p class="content">${post.Content}</p>
  <div class="date">Publié le ${new Date(post.CreatedAt).toLocaleDateString("fr-FR")}</div>
  `;

  setDocumentTitle(post.Title);
  document.querySelector(".author").addEventListener("click", () => {
    if (post.Author === "Anonyme") return;
    window.location.href = `/front/profil/profil.html?user=${post.Author}`;
  });

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

  if (user.isLoggedIn) {
    try {
      if (user.id === post.OwnerID) {
        const deleteBtn = document.createElement("button");
        deleteBtn.textContent = "Supprimer le post";
        deleteBtn.classList.add("delete-button");

        deleteBtn.addEventListener("click", async () => {
          const confirmation = confirm("Es-tu sûr de vouloir supprimer ce post ?");
          if (!confirmation) return;

          try {
            const response = await fetch(`/api/delete-post?id=${post.ID}`, {
              method: "DELETE"
            });

            if (response.ok) {
              window.location.href = "/front/post-list/postlist.html";
            } else {
              alert("Erreur lors de la suppression du post.");
            }
          } catch (error) {
            console.error("Erreur réseau :", error);
            alert("Erreur réseau lors de la suppression.");
          }
        });

        postContainer.appendChild(deleteBtn);

        const editBtn = document.createElement("button");
        editBtn.textContent = "✏️";
        editBtn.id = ("edit-button");
        document.getElementById("edit-post-button-container").appendChild(editBtn);

        const modal = document.getElementById("edit-post-modal");
        const closeModal = modal.querySelector(".close-edit-button");

        editBtn.addEventListener("click", () => {
          document.getElementById("edit-title").value = post.Title;
          document.getElementById("edit-content").value = post.Content;
          document.getElementById("edit-themes").value = post.Themes;
          modal.classList.remove("hidden");
        });

        closeModal.addEventListener("click", () => {
          modal.classList.add("hidden");
        });

        const editForm = document.getElementById("edit-post-form");
        editForm.addEventListener("submit", async (e) => {
          e.preventDefault();

          const updatedPost = {
            id: post.ID,
            title: editForm.title.value,
            content: editForm.content.value,
            theme: editForm.themes.value
          };

          try {
            const res = await fetch("/api/modify-post", {
              method: "PUT",
              headers: {
                "Content-Type": "application/json"
              },
              body: JSON.stringify(updatedPost)
            });

            if (res.ok) {
              modal.classList.add("hidden");
              location.reload();
            } else {
              alert("Erreur lors de la modification du post.");
            }
          } catch (err) {
            console.error("Erreur réseau :", err);
            alert("Erreur réseau.");
          }
        });
      }
    } catch (error) {
      console.error("Erreur lors de la récupération du profil :", error);
      alert("Erreur lors de la récupération du profil.");
    }
  } else {
    const likeBtn = document.querySelector('.like-btn');
    likeBtn.style.display = "none"; // Cacher le bouton de like si l'utilisateur n'est pas connecté
  }
  const likeSpan = document.createElement('span');
  likeBtn.appendChild(likeSpan);
  updateLikeCount(post.ID);
}

// load post & comments
document.addEventListener("DOMContentLoaded", async () => {
  if (!user.isLoggedIn || user.name === null)
    await user.getProfile();

  const postContainer = document.getElementById("dynamic-post");

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
    renderPostDetail(post);

    const likeBtn = document.querySelector('.like-btn');
    let likeSpan = document.createElement('span');
    likeSpan.id = `like-count-${post.ID}`;
    likeBtn.appendChild(likeSpan);

  } catch (err) {
    console.error("Erreur:", err);
    postContainer.innerHTML = "<p>Erreur lors du chargement du post.</p>";
  }
  await displayComments();
});

// load header buttons
document.addEventListener("DOMContentLoaded", async () => {
  const authContainer = document.querySelector(".auth-buttons");

  if (!user.isLoggedIn || user.name === null)
    await user.getProfile();

  const buttons = [
    { text: "Connexion", url: "/front/login/login.html", id: "connexion-button" },
    { text: "Inscription", url: "/front/register/register.html", id: "inscription-button" },
    { text: "Profil", url: `/front/profil/profil.html?user=${user.name}`, id: "profil-button" },
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


document.addEventListener("DOMContentLoaded", () => {
  const addCommentBtn = document.querySelector(".add-comment-btn");
  const commentForm = document.getElementById("comment-form");

  addCommentBtn.addEventListener("click", () => {
    commentForm.classList.toggle("hidden");
  });

  commentForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const content = document.getElementById("comment-content").value;

    if (!content.trim()) {
      alert("Le contenu du commentaire ne peut pas être vide.");
      return;
    }

    try {
      const res = await fetch(`/api/add-comment`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({ postId, content })
      });

      if (res.ok) {
        commentForm.reset();
        commentForm.classList.add("hidden");

        await displayComments();
      } else {
        alert("Erreur lors de l'ajout du commentaire.");
      }
    } catch (err) {
      console.error("Erreur réseau :", err);
      alert("Erreur réseau lors de l'ajout du commentaire.");
    }
  });
});

async function displayComments() {
  const commentsContainer = document.getElementById("comments-container");

  if (!commentsContainer || !postId) {
    console.error("Élément #comments-container introuvable ou ID de post manquant");
    return;
  }

  try {
    let res = await fetch(`/api/comments/${postId}`);
    if (!res.ok) {
      return;
    }
    let comments = await res.json();
    if (comments === null || comments.length === 0) {
      commentsContainer.innerHTML = "<p>Aucun commentaire pour ce post.</p>";
      return;
    }

    if (!user.isLoggedIn)
      document.querySelector(".add-comment-btn").style.display = "none";

    commentsContainer.innerHTML = "";
    comments.forEach(comment => {
      const commentDiv = document.createElement("div");
      commentDiv.className = "comment-block";
      commentDiv.innerHTML = `
          <img class="avatar" src="/api/get_avatar/${comment.Author}">
          <div class="comment-author">${comment.Author || "Auteur inconnu"}</div>
          <div class="comment-content">${comment.Content}</div>
          <div class="comment-date">Publié le ${new Date(comment.CreatedAt).toLocaleDateString("fr-FR")}</div>
        `;
      commentDiv.querySelector(".comment-author").addEventListener("click", (event) => {
        if (comment.Author === "Anonyme") return;
        window.location.href = `/front/profil/profil.html?user=${comment.Author}`;
      });

      if (user.isLoggedIn && user.id === comment.OwnerID) {
        const deleteBtn = document.createElement("button");
        deleteBtn.textContent = "Supprimer";
        deleteBtn.className = "delete-comment-button";

        deleteBtn.addEventListener("click", async () => {
          const confirmation = confirm("Es-tu sûr de vouloir supprimer ce commentaire ?");
          if (!confirmation) return;

          try {
            const response = await fetch(`/api/delete-comment?id=${comment.ID}`, {
              method: "DELETE"
            });

            if (response.ok) {
              await displayComments();
            } else {
              alert("Erreur lors de la suppression du commentaire.");
            }
          } catch (error) {
            console.error("Erreur réseau :", error);
            alert("Erreur réseau lors de la suppression.");
          }
        });

        const editBtn = document.createElement("button");
        editBtn.textContent = "Modifier";
        editBtn.className = "edit-comment-button";

        editBtn.addEventListener("click", async () => {
          const newContent = prompt("Modifier le commentaire :", comment.Content);
          if (newContent === null || newContent.trim() === "") {
            alert("Le contenu du commentaire ne peut pas être vide.");
            return;
          }

          try {
            const response = await fetch(`/api/edit-comment`, {
              method: "PUT",
              headers: {
                "Content-Type": "application/json"
              },
              body: JSON.stringify({ id: comment.ID, content: newContent })
            });

            if (response.ok) {
              await displayComments();
            } else {
              alert("Erreur lors de la modification du commentaire.");
            }
          } catch (error) {
            console.error("Erreur réseau :", error);
            alert("Erreur réseau lors de la modification.");
          }
        });

        commentDiv.appendChild(deleteBtn);
        commentDiv.appendChild(editBtn);
      }

      commentsContainer.appendChild(commentDiv);
    });
  } catch (err) {
    console.error("BBBErreur lors de la récupération des commentaires :", err);
    commentsContainer.innerHTML = "<p>Erreur lors du chargement des commentaires.</p>";
  }
}