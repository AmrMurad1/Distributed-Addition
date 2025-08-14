package com.distributedaddition.serviceb.service;
import lombok.AllArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@AllArgsConstructor
@RestController
@RequestMapping("/numbers")
public class FileController {

    private final FileStorage fileStorage;

    @PostMapping("/add")
    public ResponseEntity<String> addNumber(@RequestParam int value){
        fileStorage.addNumber(value);
        return ResponseEntity.ok("Number added successfully");
    }

    @GetMapping("/last")
    public ResponseEntity<Integer> getLastNumber(){
        int lastNumber = fileStorage.getLastNumber();
        return ResponseEntity.ok(lastNumber);
    }

}
