document.getElementById('login-form').addEventListener('submit', async (event) => {
    event.preventDefault(); 

    const email = document.getElementById('email').value.trim();
    const password = document.getElementById('password').value;

    const response = await fetch('/front/login/login.html', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include',
        body: JSON.stringify({ email, password })
    });

    if (response.ok) {
        window.location.href = '/';
    } else {
        const errorMessage = await response.text();
        alert("Erreur : " + errorMessage);
    }
});
