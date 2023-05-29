package github.com.oxxi.bechmark.controllers;


import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.UUID;

@RequestMapping("/")
@RestController
@CrossOrigin(origins = "*")
public class SecondServiceController {



    @GetMapping()
    public ResponseEntity Get() {
        UUID uuid = UUID.randomUUID();
        return ResponseEntity.ok(uuid.toString());
    }

}
