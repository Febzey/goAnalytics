(() => {
    "use strict";

    interface UserData {
        userAgent: string;
        url: string;
        referrer: string;
    }

    interface AnalyticsData {
        event: string;
        userData: UserData;
    }

    /**
     * Encoding data to send in img get.
     */
    const encodeData = (data: Record<string, string | number | boolean>): string => {
        return Object.entries(data)
            .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`)
            .join('&');
    };
    

    /**
     * Sending our data in an image get request to backend server.
     * @param data 
     */
    const sendAnalyticsData = (data: AnalyticsData): void => {
        const img = new Image();
        img.src = `http://localhost:8080/analytics?${encodeData({
            event: data.event,
            userData: JSON.stringify(data.userData),
        })}`;
    };

    /**
     * Getting some data on the client viewing the page.
     * @returns 
     */
    const gatherUserData = (): UserData => {
        return {
            userAgent: navigator.userAgent,
            url: window.location.href,
            referrer: document.referrer,
        };
    };

    /**
     * Fires when the page loads.
     * Sends data to server.
     */
    const handlePageLoad = (): void => {
        const userData = gatherUserData();
        const analyticsData: AnalyticsData = {
            event: 'pageLoad',
            userData,
        };
        sendAnalyticsData(analyticsData);
    };

    /**
     * Fires when the route is changed
     * Sends data to server.
     */
    const handleRouteChange = (): void => {
        const userData = gatherUserData();
        const analyticsData: AnalyticsData = {
            event: 'routeChange',
            userData,
        };
        sendAnalyticsData(analyticsData);
    };

    /**
     * Setting page change listener for Singe Page Applications such as
     * react, vue, angular etc.
     */
    const setupSPAListener = (): void => {
        const stateListener = (type: string) => {
            const orig = history[type as keyof History];
            return function (this: Window & typeof globalThis) {
                const rv = orig.apply(this, arguments as any);
                const event = new Event(type);
                (event as any).arguments = arguments;
                dispatchEvent(event);
                return rv;
            };
        };
        
        
        history.pushState = stateListener('pushState');
        window.addEventListener('pushState', handleRouteChange);

        if ('onhashchange' in window) {
            window.onhashchange = handleRouteChange;
        }
    };

    /**
     * Listen for when the page loads for first time.
     */
    window.addEventListener('load', handlePageLoad);

    // Check if it's an SPA and set up appropriate listeners
    if (history.pushState !== undefined) {
        setupSPAListener();
    }
    
})();
