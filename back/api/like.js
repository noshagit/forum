export { updateLikeCount, like, unlike, likeStatus };

function getCookie(name) {
    let value = "; " + document.cookie;
    let parts = value.split("; " + name + "=");
    if (parts.length === 2) {
        return parts.pop().split(";").shift();
    }
    return null;
}

async function getUserId() {
    const sessionToken = getCookie("session_token");
    if (!sessionToken) {
        return null;
    }

    try {
        const response = await fetch("/get-profile", {
            method: "GET",
            credentials: "include"
        });
        const data = await response.json();

        if (!data.success) {
            alert("Erreur lors de la récupération des informations du profil.");
            return null;
        } else {
            return data.profile.id;
        }
    } catch (error) {
        console.error("Erreur:", error);
        alert("Erreur lors de la récupération des données.");
        return null;
    }
}

async function updateLikeCount(postId) {
    try {
        const res = await fetch(`/posts/${postId}/like_count`);
        const data = await res.json();
        document.getElementById(`like-count-${postId}`).textContent = data.likes;
    } catch (err) {
        console.error('Failed to fetch like count:', err);
    }
}

async function likeStatus(postId) {
    const userId = await getUserId();
    if (!userId) {
        console.error('User not authenticated');
        return;
    }
    try {
        const res = await fetch(`/posts/is_liked`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ userId: userId, postId: postId }),
            credentials: 'include'
        });
        const data = await res.json();
        return data.liked;
    } catch (err) {
        console.error('Failed to check like status:', err);
        return false;
    }
}

async function like(postId) {
    const userId = await getUserId();
    if (!userId) {
        console.error('User not authenticated');
        return;
    }
    try {
        await fetch(`/posts/like`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ userId: userId, postId: postId }),
            credentials: 'include'
        });
        await updateLikeCount(postId);
    } catch (err) {
        console.error('Failed to like:', err);
    }
}

async function unlike(postId) {
    const userId = await getUserId();
    if (!userId) {
        console.error('User not authenticated');
        return;
    }
    try {
        await fetch(`/posts/unlike`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ userId: userId, postId: postId }),
            credentials: 'include'
        });
        await updateLikeCount(postId);
    } catch (err) {
        console.error('Failed to unlike:', err);
    }
}