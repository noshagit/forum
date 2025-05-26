function getCookie(name) {
    let value = "; " + document.cookie;
    let parts = value.split("; " + name + "=");
    if (parts.length === 2) {
        return parts.pop().split(";").shift();
    }
    return null;
}

function loadProfile() {
    const sessionToken = getCookie("session_token");

    if (sessionToken) {
        fetch("/get-profile", {
            method: "GET",
            credentials: "include"
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                document.getElementById("username").value = data.profile.username;
                document.getElementById("email").value = data.profile.email;
            } else {
                alert("Erreur lors de la récupération des informations du profil.");
            }
        })

        .catch(error => {
            console.error("Erreur:", error);
            alert("Erreur lors de la récupération des données.");
        });
    } else {
        alert("Vous n'êtes pas connecté.");
        window.location.href = "/front/login/login.html";
    }
}

function logout() {
    fetch("/logout", {
        method: "POST"
    })
    .then(() => {
        alert("Déconnexion réussie !");
        window.location.href = "/";
    })
    .catch(error => {
        console.error("Erreur de déconnexion:", error);
        alert("Erreur lors de la déconnexion.");
    });
}

document.addEventListener("DOMContentLoaded", loadProfile);

function deleteProfile() {
    if (!confirm("Êtes-vous sûr de vouloir supprimer votre compte ? Cette action est irréversible.")) {
        return;
    }

    fetch("/delete-profile", {
        method: "POST",
        credentials: "include"
    })
    .then(response => {
        if (response.ok) {
            alert("Compte supprimé avec succès !");
            window.location.href = "/"; 
        } else {
            alert("Erreur lors de la suppression du compte.");
        }
    })
    .catch(error => {
        console.error("Erreur lors de la suppression :", error);
        alert("Erreur lors de la suppression du compte.");
    });
}

function updateProfile() {
    const username = document.getElementById("username").value.trim();
    const email = document.getElementById("email").value.trim();

    if (!username || !email) {
        alert("Veuillez remplir tous les champs.");
        return;
    }

    if (!email.includes("@")) {
        alert("Veuillez entrer une adresse e-mail valide.");
        return;
    }

    fetch("/update-profile", {
        method: "POST",
        credentials: "include",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ username, email })
    })
    .then(response => {
        if (response.ok) {
            alert("Profil mis à jour avec succès !");
        } else {
            return response.text().then(text => { throw new Error(text) });
        }
    })
    .catch(error => {
        alert("Erreur lors de la mise à jour : " + error.message);
    });
}

document.querySelector(".cta").addEventListener("click", updateProfile);
