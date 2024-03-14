interface ClientMetaData {
    userAgent: string;
    url: string;
    referrer: string;
    device_width: number;
    device_height: number;
}

interface UnloadEventData {
}

interface LoadEventData {

}

interface ButtonClickData {

}

type eventTypes = "load" | "unload" | "pushstate" | "onhashchange";
type eventDataTypes = UnloadEventData | LoadEventData | ButtonClickData

interface AnalyticsPayload {
    // The type of analytic event
    event: eventTypes

    // Data about the client viewing the page.
    client_meta_data: ClientMetaData;

    // Data for the certain event.
    event_data: eventDataTypes
}

/**
 * Encoding data to send in img get.
 */
const encodeData = (data: any): string => {
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
    const gatherClientMetaData = async (): Promise<ClientMetaData> => {
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
    const sendAnalyticsData = async (event: string, event_data: eventDataTypes) => {

        const clientMeta = await gatherClientMetaData();

        const img = new Image();

        const payload = JSON.stringify({
            event: event,
            client_meta_data: clientMeta,
            event_data: event_data,
        } as AnalyticsPayload)

        const queryString = `http://localhost:8080/analytics?${encodeData({
            analytics_payload: payload
        })}`;

        img.src = queryString;
    };

    /**
     * Fires when the route is changed or page loads or hash change
     * Sends data to server.
     */
    const handleRouteChange = (type: eventTypes) => {
        sendAnalyticsData(type, {});
    };

    const handleUnloadPage = (type: string) => {

        console.log("Unloaded page bruh")

        const pageUnloadData: UnloadEventData = {

        }
        sendAnalyticsData(type, pageUnloadData)
    }

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
