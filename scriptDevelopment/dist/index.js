"use strict";
(() => {
    "use strict";
    /**
     * Encoding data to send in img get.
     */
    const encodeData = (data) => {
        return Object.entries(data)
            .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`)
            .join('&');
    };
    /**
     * Sending our data in an image get request to backend server.
     * @param data
     */
    const sendAnalyticsData = (data) => {
        const img = new Image();
        img.src = `http://localhost:8080/analytics?${encodeData({
            data: JSON.stringify(data.EventData),
        })}`;
    };
    /**
     * Getting some data on the client viewing the page.
     * @returns
     */
    const gatherEventData = async (eventType) => {
        return {
            event: eventType,
            userAgent: navigator.userAgent,
            url: window.location.href,
            referrer: document.referrer,
        };
    };
    /**
     * Fires when the route is changed or page loads or hash change
     * Sends data to server.
     */
    const handleRouteChange = (type) => {
        gatherEventData(type)
            .then((eventData) => {
            const analyticsData = {
                event: 'routeChange',
                EventData: eventData,
            };
            sendAnalyticsData(analyticsData);
        })
            .catch((error) => {
            console.error('Error handling route change:', error);
        });
    };
    /**
     * Setting page change listener for Singe Page Applications such as
     * react, vue, angular etc.
     */
    const setupSPAListener = () => {
        const stateListener = (type) => {
            const orig = history[type];
            return function () {
                const rv = orig.apply(this, arguments);
                const event = new Event(type);
                event.arguments = arguments;
                dispatchEvent(event);
                return rv;
            };
        };
        history.pushState = stateListener('pushState');
        window.addEventListener('pushState', () => handleRouteChange("pushstate"));
        if ('onhashchange' in window) {
            window.onhashchange = () => handleRouteChange("onhashchange");
        }
    };
    /**
     * Listen for when the page loads for first time.
     */
    window.addEventListener('load', () => handleRouteChange("load"));
    // Check if it's an SPA and set up appropriate listeners
    if (history.pushState !== undefined) {
        setupSPAListener();
    }
})();
