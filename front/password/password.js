document.addEventListener("DOMContentLoaded", async () => {

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
    const button = document.querySelector(".cta");

    button.addEventListener("click", async () => {
        const password = document.querySelectorAll("input")[0].value;
        const confirmPassword = document.querySelectorAll("input")[1].value;

        if (password !== confirmPassword) {
            alert("Les mots de passe ne correspondent pas.");
            return;
        }

        const response = await fetch("/api/change-password", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ password })
        });

        const data = await response.json();

        if (response.ok) {
            window.location.href = `/front/profil/profil.html?user=${username}`;
        } else {
            alert(`Erreur : ${data.error}`);
        }
    });
});
