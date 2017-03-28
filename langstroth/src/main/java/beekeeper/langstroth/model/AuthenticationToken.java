package beekeeper.langstroth.model;

/**
 * @author willstuckey
 * @date 3/27/17
 * <p></p>
 */
public interface AuthenticationToken {
    String getAuthToken();
    String getAuthTokenExpiration();
}
