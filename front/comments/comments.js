import { like, unlike, updateLikeCount, likeStatus } from "/api/like";

const urlParams = new URLSearchParams(window.location.search);
const postId = urlParams.get('id');

function setDocumentTitle(title) {
  if (title) {
    document.title = title;
  }
}

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

  const isLoggedIn = document.cookie.includes("session_token=");
  if (isLoggedIn) {
    try {
      const profileRes = await fetch("/get-profile");
      if (!profileRes.ok) throw new Error("Erreur lors de la récupération du profil");

      const user = await profileRes.json();

      if (user.ID === post.ID_Owner) {
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
              alert("Post supprimé avec succès.");
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
        editBtn.textContent = "Modifier le post";
        editBtn.classList.add("edit-button");
        document.getElementById("edit-post-button-container").appendChild(editBtn);

        // Prépare la modal
        const modal = document.getElementById("edit-post-modal");
        const closeModal = modal.querySelector(".close-edit-button");

        // Bouton ouvrir la pop-up
        editBtn.addEventListener("click", () => {
          document.getElementById("edit-title").value = post.Title;
          document.getElementById("edit-content").value = post.Content;
          document.getElementById("edit-themes").value = post.Theme || "";
          modal.classList.remove("hidden");
        });

        // Bouton fermer la pop-up
        closeModal.addEventListener("click", () => {
          modal.classList.add("hidden");
        });

        // Soumission du formulaire de modification
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
              alert("Post modifié avec succès !");
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
      } else {
        const likeBtn = document.querySelector('.like-btn');
        likeBtn.style.display = "none"; // Cacher le bouton de like si l'utilisateur n'est pas l'auteur
      }
    } catch (error) {
      console.error("Erreur lors de la récupération du profil :", error);
      alert("Erreur lors de la récupération du profil.");
    }
  } else {
    const likeBtn = document.querySelector('.like-btn');
    likeBtn.style.display = "none"; // Cacher le bouton de like si l'utilisateur n'est pas connecté
  }
  updateLikeCount(post.ID);
  const likeSpan = document.createElement('span');
  likeBtn.appendChild(likeSpan);
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
      <div class="date">Publié le ${new Date(post.CreatedAt).toLocaleDateString("fr-FR")}</div>
    `;
          console.log("Post chargé avec succès :", post);
          updateLikeCount(post.ID);

          const likeBtn = document.querySelector('.like-btn');
          let likeSpan = document.createElement('span');
          likeSpan.id = `like-count-${post.ID}`;
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