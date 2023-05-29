package github.com.oxxi.bechmark.controllers;


import lombok.Builder;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
@Builder
public class ResponseMinio {
    private String Etag;
    private String versionId;
}
