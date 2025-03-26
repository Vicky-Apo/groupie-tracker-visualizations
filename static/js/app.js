document.addEventListener("DOMContentLoaded", () => {
    // 1. Trigger Event button code
    const btn = document.getElementById("eventBtn");
    if (btn) {
        btn.addEventListener("click", () => {
            fetch("/trigger-event", {
                method: "POST",
                headers: { "Content-Type": "application/json" }
            })
            .then(response => response.json())
            .then(data => {
                alert(data.status);
            })
            .catch(err => {
                console.error(err);
                alert("Error triggering event");
            });
        });
    }

    // 2. Load More functionality
    let offset = 10;  // We already have 10 displayed on page load
    const limit = 10; // We'll load 10 more per request
    const loadMoreBtn = document.getElementById("loadMoreBtn");
    const artistList = document.getElementById("artistList");

    if (loadMoreBtn) {
        loadMoreBtn.addEventListener("click", () => {
            console.log("Load More button clicked. Offset:", offset);
            fetch(`/api/artists?offset=${offset}&limit=${limit}`)
                .then(response => {
                    console.log("Response status:", response.status);
                    return response.json();
                })
                .then(data => {
                    console.log("Data received:", data);
                    // If no data returned, hide the button.
                    if (data.length === 0) {
                        loadMoreBtn.style.display = "none";
                        console.log("No more artists to load. Hiding button.");
                        return;
                    }

                    // For each new artist, create a card.
                    data.forEach(artist => {
                        const card = document.createElement("div");
                        card.className = "artist-card";

                        // Convert spaces to underscores in the artist's name
                        const artistUrl = "/artist/" + artist.name.replace(/\s+/g, "-");

                        card.innerHTML = `
                            <a href="${artistUrl}">
                                <img src="${artist.image}" alt="${artist.name}" class="artist-img" />
                                <h3>${artist.name}</h3>
                            </a>
                        `;
                        artistList.appendChild(card);
                    });

                    // Increase offset by limit for next fetch
                    offset += limit;
                    console.log("New offset:", offset);
                })
                .catch(err => {
                    console.error("Error loading more artists:", err);
                });
        });
    }

    // 3. Existing toggle code (optional)
    document.querySelectorAll('.toggle-section').forEach(function(button) {
        button.addEventListener('click', function() {
            var targetId = this.getAttribute('data-target');
            var section = document.getElementById(targetId);
            if (section.style.display === 'none') {
                section.style.display = 'block';
                this.innerText = 'Hide';
            } else {
                section.style.display = 'none';
                this.innerText = 'Show';
            }
        });
    });
});
