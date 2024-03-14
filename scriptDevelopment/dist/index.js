"use strict";
/**
 * Encoding data to send in img get.
 */
const encodeData = (data) => {
    return Object.entries(data)
        .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`)
        .join('&');
};
(() => {
    //Ideally we should create some sort of cache here, 
    //Figure out a way to handle page view durations
    //When a user "loads" a page, we should start the duration counter
    //when a user "unloads" we should end the duration counter and send it.
    //We will send the duration amount in the unload event payload.
    /**
     * Getting some data on the client viewing the page.
     * @returns
     */
    const gatherClientMetaData = async () => {
        return {
            userAgent: navigator.userAgent,
            url: window.location.href,
            referrer: document.referrer,
            device_height: window.innerHeight,
            device_width: window.innerWidth
        };
    };
    /**
     * Sending our data in an image get request to backend server.
     * @param data
     */
    const sendAnalyticsData = async (event, event_data) => {
        const clientMeta = await gatherClientMetaData();
        const img = new Image();
        const payload = JSON.stringify({
            event: event,
            client_meta_data: clientMeta,
            event_data: event_data,
        });
        const queryString = `http://localhost:8080/analytics?${encodeData({
            analytics_payload: payload
        })}`;
        img.src = queryString;
    };
    /**
     * Fires when the route is changed or page loads or hash change
     * Sends data to server.
     */
    const handleRouteChange = (type) => {
        sendAnalyticsData(type, {});
    };
    const handleUnloadPage = (type) => {
        console.log("Unloaded page bruh");
        const pageUnloadData = {};
        sendAnalyticsData(type, pageUnloadData);
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
    window.addEventListener('unload', () => handleUnloadPage("unload"));
    // Check if it's an SPA and set up appropriate listeners
    if (history.pushState !== undefined) {
        setupSPAListener();
    }
})();
