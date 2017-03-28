package beekeeper.langstroth.api.auth;

import org.springframework.security.core.AuthenticationException;
import org.springframework.security.web.AuthenticationEntryPoint;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;

/**
 * @author willstuckey
 * @date 3/27/17
 * <p></p>
 */
public class DefaultAuthenticationEntryPoint implements AuthenticationEntryPoint {
    private static final String AUTH_FAILED_MSG =
            "Unauthorized: Authentication token was either missing or invalid.";

    @Override
    public void commence(HttpServletRequest request,
                         HttpServletResponse response,
                         AuthenticationException authException)
            throws IOException, ServletException {
        response.sendError(HttpServletResponse.SC_UNAUTHORIZED, AUTH_FAILED_MSG);
    }
}