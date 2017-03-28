package beekeeper.langstroth.api;

/**
 * @author willstuckey
 * @date 3/27/17
 * <p></p>
 */
public final class APIEndPoints {
    public static final class Version {
        public static final String V1 = "/v1";
    }

    public static final class Modules {
        public static final String AUTH = "/auth";

        public static final class Auth {
            public static final String WITH_CREDENTIALS =  "/cred";
            public static final String WITH_TOKEN =   "/tok";
            public static final String DEAUTH =     "/deauth";
        }
    }
}
