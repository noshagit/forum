document.addEventListener("DOMContentLoaded", () => {
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
            window.location.href = "/front/profil/profil.html";
        } else {
            alert(`Erreur : ${data.error}`);
        }
    });
});
