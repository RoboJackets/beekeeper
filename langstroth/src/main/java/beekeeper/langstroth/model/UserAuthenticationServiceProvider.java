package beekeeper.langstroth.model;

/**
 * @author willstuckey
 * @date 3/27/17
 * <p>This interface defines the functionality for authentication</p>
 */
public interface UserAuthenticationServiceProvider {
    /**
     * generates a valid token for API request authentication from
     * account identification and password
     * @param id account identification, exact content defined by
     *           implementation
     * @param password account password, plaintext
     * @return valid API authentication token
     */
    AuthenticationToken getToken(final String id, final String password);

    /**
     * generates a valid token for API request authentication from
     * a previous auth token. Implementation defines what happens
     * to the old token.
     * @param token old but unexpired token
     * @return new API auth token
     */
    AuthenticationToken getToken(final String token);

    /**
     * check if an auth token is valid
     * @param token auth token
     * @return validity
     */
    boolean authenticated(final AuthenticationToken token);


}
