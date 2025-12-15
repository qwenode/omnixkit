// Simple SPA router demo
(function() {
    const routes = {
        '/': {
            title: 'Home',
            content: '<h2>Welcome to kitspa Demo</h2><p>This is a simple single-page application served by kitspa.</p><p>kitspa handles SPA routing by returning index.html for all non-asset routes, allowing client-side routing to work correctly.</p>'
        },
        '/about': {
            title: 'About',
            content: '<h2>About kitspa</h2><p>kitspa is a Go package for serving single-page applications with Gin.</p><ul><li>Supports embedded filesystem (embed.FS)</li><li>Configurable static assets path</li><li>SPA fallback routing</li><li>Blocked path prefixes for security</li></ul>'
        },
        '/contact': {
            title: 'Contact',
            content: '<h2>Contact</h2><p>This is the contact page.</p><p>Try navigating between pages - notice how the URL changes but the page does not reload.</p>'
        }
    };

    function navigate(path) {
        const route = routes[path] || {
            title: '404',
            content: '<h2>Page Not Found</h2><p>The requested page does not exist.</p>'
        };

        document.getElementById('content').innerHTML = route.content;
        document.getElementById('current-path').textContent = path;

        // Update active nav
        document.querySelectorAll('.nav a').forEach(a => {
            a.classList.toggle('active', a.getAttribute('href') === path);
        });
    }

    // Handle navigation clicks
    document.addEventListener('click', function(e) {
        if (e.target.tagName === 'A' && e.target.getAttribute('href').startsWith('/')) {
            e.preventDefault();
            const path = e.target.getAttribute('href');
            history.pushState(null, '', path);
            navigate(path);
        }
    });

    // Handle browser back/forward
    window.addEventListener('popstate', function() {
        navigate(window.location.pathname);
    });

    // Initial render
    navigate(window.location.pathname);
})();
